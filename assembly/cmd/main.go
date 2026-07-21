package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	assemblyService "github.com/Reensef/go-microservices-course/assembly/internal/service/assembly"
)

func main() {
	_ = assemblyService.NewService()

	log.Println("🚀 assembly service started")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down assembly service...")
	log.Println("✅ Service stopped")
}
