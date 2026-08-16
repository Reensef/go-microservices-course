package order

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
)

func (s *service) CancelOrder(ctx context.Context, orderUuid, requesterUuid string) error {
	if uuid.Validate(orderUuid) != nil {
		return model.ErrOrderUuidInvalidFormat
	}

	logger.Info("cancelling order",
		zap.String("order_uuid", orderUuid),
		zap.String("requester_uuid", requesterUuid),
	)

	order, err := s.orderRepo.GetOrderByUUID(ctx, orderUuid)
	if err != nil {
		return err
	}

	if order.Info.UserUuid != requesterUuid {
		return model.ErrOrderAccessDenied
	}

	if order.Info.Status == model.OrderStatus_PAID {
		return model.ErrOrderAlreadyPaid
	}

	err = s.orderRepo.CancelOrder(ctx, orderUuid)
	if err != nil {
		logger.Error("failed to cancel order", zap.String("order_uuid", orderUuid), zap.Error(err))
		return err
	}

	logger.Info("order cancelled", zap.String("order_uuid", orderUuid))
	return nil
}
