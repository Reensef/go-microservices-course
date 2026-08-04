package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	converter "github.com/Reensef/go-microservices-course/order/internal/client/grpc/iam/v1/converter"
	"github.com/Reensef/go-microservices-course/order/internal/model"
	iamGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

func (c *iamClient) Whoami(ctx context.Context, sessionUUID string) (model.User, error) {
	response, err := c.service.Whoami(ctx,
		&iamGrpc.WhoamiRequest{
			SessionUuid: sessionUUID,
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && (st.Code() == codes.NotFound || st.Code() == codes.InvalidArgument) {
			return model.User{}, model.ErrInvalidSession
		}

		return model.User{}, fmt.Errorf("iam whoami: %w", err)
	}

	return converter.ToModelUser(response.GetUser()), nil
}
