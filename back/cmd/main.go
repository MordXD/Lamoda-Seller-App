package main

import (
	"log"

	"github.com/lamoda-seller-app/internal/config"
	"github.com/lamoda-seller-app/internal/server"
)

func main() {
	cfg := config.Load()

	s, err := server.Init(cfg)
	if err != nil {
		log.Fatalf("❌ Server initialization failed: %v", err)
	}

	s.Run()
}