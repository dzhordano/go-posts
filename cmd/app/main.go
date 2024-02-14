package main

import (
	"log"

	"github.com/dzhordano/go-posts/internal/common/server"
	delivery "github.com/dzhordano/go-posts/internal/delivery/http"
)

func main() {

	handlers := delivery.NewHandler()

	srv := server.NewServer(handlers.Init())

	log.Fatal(srv.Run())
}
