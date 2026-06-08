package main

import (
	"ServidorA/gen/chat"
	"ServidorA/gen/desabafo"
	"ServidorA/internal/connection"
	"ServidorA/internal/persistence"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {

	if os.Getenv("ROOT") == "" {
		godotenv.Load() // só carrega o .env se a var ainda não estiver definida
	}

	persistence.StartDataBase()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	desabafo.RegisterFeedServer(s, &connection.FeedServer{})
	chat.RegisterChatServiceServer(s, &connection.ChatServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
