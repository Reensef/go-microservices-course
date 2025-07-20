package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/payment/v1"
)

// Реализует gRPC сервис для оплыты заказов
type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func NewPaymentService() *paymentService {
	return &paymentService{}
}

// Обработчик запроса на оплату
func (s *paymentService) PayOrder(
	context.Context,
	*paymentV1.PayOrderRequest,
) (*paymentV1.PayOrderResponse, error) {
	transaction_uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}

	log.Println("Оплата прошла успешно, transaction_uuid:", transaction_uuid.String())

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transaction_uuid.String(),
	}, nil
}
