package main

import (
	"log"

	delivery "github.com/dzhordano/go-posts/internal/delivery/http"
	"github.com/dzhordano/go-posts/pkg/server"
)

func main() {
	handlers := delivery.NewHandler()

	srv := server.NewServer(handlers.Init())

	log.Fatal(srv.Run())
}
