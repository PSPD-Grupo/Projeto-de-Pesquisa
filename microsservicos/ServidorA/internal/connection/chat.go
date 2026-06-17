package connection

import (
	"ServidorA/gen/chat"
	"log"
	"sync"

	"google.golang.org/grpc"
)

type ChatServer struct {
	chat.UnimplementedChatServiceServer
	mu    sync.Mutex
	salas map[string]*room
}

type room struct {
	id       string
	Clientes map[string]cliente
}

type cliente struct {
	nickname string
	room_id  string
	stream   grpc.BidiStreamingServer[chat.ChatMessage, chat.ChatMessage]
}

func (s *ChatServer) broadcast(roomId string, msg *chat.ChatMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.salas[roomId]
	if !ok {
		return
	}

	for _, client := range r.Clientes {
		go func(c cliente, m *chat.ChatMessage) {
			if err := c.stream.Send(m); err != nil {
				log.Printf("erro ao enviar mensagem para %s: %v", c.nickname, err)
			}
		}(client, msg)
	}
}

func (s *ChatServer) Chat(stream grpc.BidiStreamingServer[chat.ChatMessage, chat.ChatMessage]) error {
	// Recebe a primeira mensagem para identificar o cliente e a sala
	firstMsg, err := stream.Recv()
	if err != nil {
		return err
	}

	roomId := firstMsg.RoomId
	clienteId := firstMsg.SenderNickname

	s.mu.Lock()
	if s.salas == nil {
		s.salas = make(map[string]*room)
	}
	r, ok := s.salas[roomId]
	if !ok {
		r = &room{
			id:       roomId,
			Clientes: make(map[string]cliente),
		}
		s.salas[roomId] = r
	}
	r.Clientes[clienteId] = cliente{
		nickname: clienteId,
		room_id:  roomId,
		stream:   stream,
	}
	s.mu.Unlock()

	log.Printf("cliente %s entrou na sala %s", clienteId, roomId)
	s.broadcast(roomId, firstMsg)

	defer func() {
		s.mu.Lock()
		if r, ok := s.salas[roomId]; ok {
			delete(r.Clientes, clienteId)
			if len(r.Clientes) == 0 {
				delete(s.salas, roomId)
			}
		}
		s.mu.Unlock()
		log.Printf("cliente %s saiu da sala %s", clienteId, roomId)
	}()

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("stream do cliente %s encerrado: %v", clienteId, err)
			return err
		}
		log.Printf("mensagem recebida de %s na sala %s: %s", msg.SenderNickname, msg.RoomId, msg.Text)
		s.broadcast(msg.RoomId, msg)
	}
}
