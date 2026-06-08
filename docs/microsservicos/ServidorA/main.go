package ServidorA

import (
	"ServidorA/internal/persistence"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	persistence.StartDataBase()

}
