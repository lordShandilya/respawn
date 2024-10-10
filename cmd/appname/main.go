package main

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/server"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()

	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Unable to connect database: %v", err)
	}

	defer db.Close()

	srv := server.NewServer(cfg)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failure to start the server: %v", err)
		}
	}()

	//TODO: Graceful shutdown on OS interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error shutting down gracefully: %v", err)
	}
}
