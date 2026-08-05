package app

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentApi "github.com/Reensef/go-microservices-course/payment/internal/api/payment/v1"
	grpcClients "github.com/Reensef/go-microservices-course/payment/internal/client/grpc"
	iamClient "github.com/Reensef/go-microservices-course/payment/internal/client/grpc/iam/v1"
	"github.com/Reensef/go-microservices-course/payment/internal/config"
	service "github.com/Reensef/go-microservices-course/payment/internal/service"
	paymentService "github.com/Reensef/go-microservices-course/payment/internal/service/payment"
	closer "github.com/Reensef/go-microservices-course/platform/pkg/closer"
	iamGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
	paymentProtoApi "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentApi     paymentProtoApi.PaymentServiceServer
	paymentService service.PaymentService

	iamClient grpcClients.IAMClient
	iamGrpc   iamGrpc.AuthServiceClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentApi(ctx context.Context) paymentProtoApi.PaymentServiceServer {
	if d.paymentApi == nil {
		d.paymentApi = paymentApi.New(d.PaymentService(ctx))
	}

	return d.paymentApi
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.New()
	}

	return d.paymentService
}

func (d *diContainer) IAMClient(ctx context.Context) grpcClients.IAMClient {
	if d.iamClient == nil {
		d.iamClient = iamClient.New(d.IAMGrpc(ctx))
	}

	return d.iamClient
}

func (d *diContainer) IAMGrpc(ctx context.Context) iamGrpc.AuthServiceClient {
	if d.iamGrpc == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IAMClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to iam service: %v\n", err))
		}

		closer.AddNamed("IAM gRPC client", func(ctx context.Context) error {
			if err := conn.Close(); err != nil {
				return err
			}
			return nil
		})

		d.iamGrpc = iamGrpc.NewAuthServiceClient(conn)
	}

	return d.iamGrpc
}
