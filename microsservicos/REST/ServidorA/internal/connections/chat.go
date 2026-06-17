package connection

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ChatMessage equivale ao proto chat.ChatMessage
type ChatMessage struct {
	RoomID         string `json:"room_id"`
	SenderNickname string `json:"sender_nickname"`
	Text           string `json:"text"`
}

type ChatServer struct {
	salas map[string]*room
}

type room struct {
	id        string
	Clientes  map[string]*cliente
	broadcast chan *ChatMessage
}

type cliente struct {
	nickname string
	room_id  string
	conn     *websocket.Conn
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		salas: make(map[string]*room),
	}
}

func (s *ChatServer) getOrCreateRoom(roomID string) *room {
	if r, ok := s.salas[roomID]; ok {
		return r
	}
	r := &room{
		id:        roomID,
		Clientes:  make(map[string]*cliente),
		broadcast: make(chan *ChatMessage, 64),
	}
	s.salas[roomID] = r
	go r.runBroadcast()
	return r
}

// runBroadcast distribui mensagens para todos os clientes da sala
func (r *room) runBroadcast() {
	for msg := range r.broadcast {
		data, err := json.Marshal(msg)
		if err != nil {
			log.Println("erro ao serializar:", err)
			continue
		}
		for _, c := range r.Clientes {
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("erro ao enviar para %s: %v", c.nickname, err)
			}
		}
	}
}

// Chat é o handler HTTP/WebSocket — equivale ao RPC Chat(stream) do gRPC
func (s *ChatServer) Chat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade falhou:", err)
		return
	}
	defer conn.Close()

	// Primeira mensagem: identificação do cliente e da sala
	_, raw, err := conn.ReadMessage()
	if err != nil {
		log.Println("erro na primeira mensagem:", err)
		return
	}

	var firstMsg ChatMessage
	if err := json.Unmarshal(raw, &firstMsg); err != nil {
		log.Println("json inválido:", err)
		return
	}

	roomID := firstMsg.RoomID
	clienteID := firstMsg.SenderNickname

	sala := s.getOrCreateRoom(roomID)
	sala.Clientes[clienteID] = &cliente{
		nickname: clienteID,
		room_id:  roomID,
		conn:     conn,
	}

	log.Printf("cliente %s entrou na sala %s", clienteID, roomID)

	// Envia a primeira mensagem para o broadcast da sala
	sala.broadcast <- &firstMsg

	defer func() {
		delete(sala.Clientes, clienteID)
		log.Printf("cliente %s saiu da sala %s", clienteID, roomID)
	}()

	// Loop de recebimento — equivale ao goroutine de Recv() do gRPC
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println("stream encerrado:", err)
			return
		}

		var msg ChatMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Println("json inválido:", err)
			continue
		}

		log.Printf("mensagem recebida de %s: %s", msg.SenderNickname, msg.Text)
		s.salas[msg.RoomID].broadcast <- &msg
	}
}


