package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/Reensef/go-microservices-course/order/internal/api/order/v1"
	inventoryServiceV1 "github.com/Reensef/go-microservices-course/order/internal/client/grpc/inventory/v1"
	paymentServiceV1 "github.com/Reensef/go-microservices-course/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/Reensef/go-microservices-course/order/internal/repository/order"
	orderService "github.com/Reensef/go-microservices-course/order/internal/service/order"
	orderV1 "github.com/Reensef/go-microservices-course/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
	"github.com/Reensef/go-microservices-course/shared/pkg/utils"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout       = 5 * time.Second
	shutdownTimeout         = 10 * time.Second
	inventoryServiceAddress = "localhost:50051"
	paymenyServiceAddress   = "localhost:50052"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	postgresURI := os.Getenv("POSTGRES_URI")

	pool, err := pgxpool.New(ctx, postgresURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer pool.Close()

	// // Создаем соединение с базой данных
	// con, err := pgx.Connect(ctx, postgresURI)
	// if err != nil {
	// 	log.Printf("failed to connect to database: %v\n", err)
	// 	return
	// }
	// defer func() {
	// 	cerr := con.Close(ctx)
	// 	if cerr != nil {
	// 		log.Printf("failed to close connection: %v\n", cerr)
	// 	}
	// }()

	// Проверяем, что соединение с базой установлено
	err = pool.Ping(ctx)
	if err != nil {
		log.Printf("database unavailable: %v\n", err)
		return
	}

	// Инициализируем мигратор
	migrationsDir := os.Getenv("ORDER_MIGRATIONS_DIR")
	pool.Config().Copy()
	migrator := utils.NewSqlMigrator(stdlib.OpenDB(*pool.Config().ConnConfig.Copy()), migrationsDir)

	err = migrator.Up()
	if err != nil {
		log.Printf("Error migrate database: %v\n", err)
		return
	}

	inventoryServiceConn, err := grpc.NewClient(
		inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryServiceConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v\n", cerr)
		}
	}()

	inventoryService := inventoryServiceV1.NewClient(
		inventoryV1.NewInventoryServiceClient(inventoryServiceConn),
	)

	paymentServiceConn, err := grpc.NewClient(
		paymenyServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentServiceConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v\n", cerr)
		}
	}()

	paymentService := paymentServiceV1.NewClient(
		paymentV1.NewPaymentServiceClient(paymentServiceConn),
	)

	repo := orderRepo.NewRepository(pool)
	service := orderService.NewService(repo, inventoryService, paymentService)
	api := orderV1API.NewAPI(service)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("Error creating OpenAPI server: %v\n", err)
		return
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	go func() {
		log.Printf("🚀 HTTP server running on port %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Server startup error: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Error while stopping the server: %v\n", err)
	}

	log.Println("✅ Server stopped")
}
