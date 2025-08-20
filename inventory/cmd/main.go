package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryApi "github.com/Reensef/go-microservices-course/inventory/internal/api/inventory/v1"
	partRepository "github.com/Reensef/go-microservices-course/inventory/internal/repository/part"
	inventoryService "github.com/Reensef/go-microservices-course/inventory/internal/service/inventory"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	mongoURI := os.Getenv("MONGO_URI")

	log.Printf("Connecting to %s...\n", mongoURI)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("failed to connect to mongo: %v\n", err)
		return
	}
	defer func() {
		cerr := mongoClient.Disconnect(ctx)
		if cerr != nil {
			log.Printf("failed to disconnect from mongo: %v\n", cerr)
		}
	}()

	// Проверяем соединение с базой данных
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}

	// Получаем базу данных
	mongoDB := mongoClient.Database("inventory")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()

	// Включаем рефлексию для отладки
	reflection.Register(s)

	repo := partRepository.NewRepository(mongoDB)
	service := inventoryService.NewService(repo)
	err = service.GenData()
	if err != nil {
		log.Printf("failed to generate data: %v\n", err)
		return
	}
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
