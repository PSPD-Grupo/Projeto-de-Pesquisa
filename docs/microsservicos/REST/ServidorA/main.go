package main

import (
	"ServidorA/internal/connections"
	"ServidorA/internal/persistence"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("ROOT") == "" {
		godotenv.Load() // só carrega o .env se a var ainda não estiver definida
	}

	persistence.StartDataBase()

	chatServer := connection.NewChatServer()
	feedServer := &connection.FeedServer{}

	http.HandleFunc("/chat", chatServer.Chat)
	http.HandleFunc("/postDesabafo", feedServer.PostDesabafo)
	http.HandleFunc("/feed", feedServer.GetFeed)

	http.ListenAndServe(":6789", nil)
}
