package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	converter "github.com/Reensef/go-microservices-course/inventory/internal/client/grpc/iam/v1/converter"
	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	iamGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

func (c *iamClient) Whoami(ctx context.Context, sessionUUID string) (model.User, error) {
	response, err := c.service.Whoami(ctx, &iamGrpc.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		status, ok := grpcstatus.FromError(err)
		if ok && (status.Code() == codes.NotFound || status.Code() == codes.InvalidArgument) {
			return model.User{}, model.ErrInvalidSession
		}

		return model.User{}, fmt.Errorf("iam whoami: %w", err)
	}

	return converter.ToModelUser(response.GetUser()), nil
}
