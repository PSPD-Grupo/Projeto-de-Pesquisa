package connection

import (
	"ServidorA/gen/chat"
	"log"

	"google.golang.org/grpc"
)

type ChatServer struct {
	chat.UnimplementedChatServiceServer
	salas map[string]room
}

type room struct {
	id        string
	Clientes  map[string]cliente
	broadcast chan *chat.ChatMessage
}

type cliente struct {
	nickname string
	room_id  string
	stream   grpc.BidiStreamingServer[chat.ChatMessage, chat.ChatMessage]
}

func (s *ChatServer) Chat(stream grpc.BidiStreamingServer[chat.ChatMessage, chat.ChatMessage]) error {
	// Canal interno para erros do goroutine de recebimento
	errCh := make(chan error, 1)

	firstMsg, err := stream.Recv()
	if err != nil {
		return err
	}

	roomId := firstMsg.RoomId
	clienteId := firstMsg.SenderNickname

	// Registra o cliente na sala
	s.salas[roomId].Clientes[clienteId] = cliente{
		nickname: firstMsg.SenderNickname,
		room_id:  roomId,
		stream:   stream,
	}
	broadcast := s.salas[roomId].broadcast

	log.Printf("cliente %s entrou na sala %s", clienteId, roomId)

	// Processa a primeira mensagem normalmente
	broadcast <- firstMsg

	// ── Goroutine de RECEBIMENTO (cliente → servidor) ──
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				errCh <- err
				return
			}

			log.Printf("mensagem recebida de %s: %s", msg.SenderNickname, msg.Text)
			s.salas[msg.RoomId].broadcast <- msg
		}
	}()

	defer func() {
		// Remove o cliente da sala ao desconectar
		delete(s.salas[roomId].Clientes, clienteId)
		log.Printf("cliente %s saiu da sala %s", clienteId, roomId)
	}()

	// ── Loop de ENVIO (servidor → cliente) ──
	for {
		select {
		case err = <-errCh:
			// Cliente desconectou ou erro no recebimento
			log.Println("stream encerrado:", err)
			return err

		case msg := <-broadcast:
			if err = stream.Send(msg); err != nil {
				log.Println("erro ao enviar:", err)
				return err
			}
		}
	}
}
