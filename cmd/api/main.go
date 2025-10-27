package main

import (
	"log"

	"github.com/thebiatriz/go-db-api/internal/server"
)

func main() {
	server, err := server.NewServer()

	if err != nil {
		log.Fatalf("Erro ao inicializar o servidor: %v", err)
	}

	server.Run(":8080")
}
