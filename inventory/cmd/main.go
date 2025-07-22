package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/Reensef/go-microservices-course/inventory/internal/api/inventory/v1"
	partRepository "github.com/Reensef/go-microservices-course/inventory/internal/repository/part"
	inventoryService "github.com/Reensef/go-microservices-course/inventory/internal/service/inventory"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()

	// Включаем рефлексию для отладки
	reflection.Register(s)

	repo := partRepository.NewRepository()
	service := inventoryService.NewService(repo)
	api := inventoryApi.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
