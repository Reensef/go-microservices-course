package grpc

import (
	"context"

	"github.com/Reensef/go-microservices-course/payment/internal/model"
)

type IAMClient interface {
	Whoami(ctx context.Context, sessionUUID string) (model.User, error)
}
