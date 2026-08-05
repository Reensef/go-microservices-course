package interceptor

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	grpcMocks "github.com/Reensef/go-microservices-course/inventory/internal/client/grpc/mocks"
	"github.com/Reensef/go-microservices-course/inventory/internal/model"
)

func noopHandler(gotCtx *context.Context) grpc.UnaryHandler {
	return func(ctx context.Context, req any) (any, error) {
		if gotCtx != nil {
			*gotCtx = ctx
		}
		return "ok", nil
	}
}

func TestAuthInterceptor_exemptMethod(t *testing.T) {
	iamClient := grpcMocks.NewMockIAMClient(t)
	interceptor := NewAuthInterceptor(iamClient)

	info := &grpc.UnaryServerInfo{FullMethod: "/grpc.health.v1.Health/Check"}

	resp, err := interceptor(context.Background(), nil, info, noopHandler(nil))

	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)
	assert.Empty(t, iamClient.Calls)
}

func TestAuthInterceptor_missingMetadata(t *testing.T) {
	iamClient := grpcMocks.NewMockIAMClient(t)
	interceptor := NewAuthInterceptor(iamClient)

	info := &grpc.UnaryServerInfo{FullMethod: "/inventory.v1.InventoryService/GetPart"}

	_, err := interceptor(context.Background(), nil, info, noopHandler(nil))

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestAuthInterceptor_malformedAuthorization(t *testing.T) {
	iamClient := grpcMocks.NewMockIAMClient(t)
	interceptor := NewAuthInterceptor(iamClient)

	info := &grpc.UnaryServerInfo{FullMethod: "/inventory.v1.InventoryService/GetPart"}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Token abc"))

	_, err := interceptor(ctx, nil, info, noopHandler(nil))

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Empty(t, iamClient.Calls)
}

func TestAuthInterceptor_invalidSession(t *testing.T) {
	iamClient := grpcMocks.NewMockIAMClient(t)
	interceptor := NewAuthInterceptor(iamClient)

	token := "session-token"
	iamClient.EXPECT().Whoami(mock.Anything, token).Return(model.User{}, model.ErrInvalidSession).Once()

	info := &grpc.UnaryServerInfo{FullMethod: "/inventory.v1.InventoryService/GetPart"}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))

	_, err := interceptor(ctx, nil, info, noopHandler(nil))

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestAuthInterceptor_iamError(t *testing.T) {
	iamClient := grpcMocks.NewMockIAMClient(t)
	interceptor := NewAuthInterceptor(iamClient)

	token := "session-token"
	iamClient.EXPECT().Whoami(mock.Anything, token).Return(model.User{}, fmt.Errorf("unavailable")).Once()

	info := &grpc.UnaryServerInfo{FullMethod: "/inventory.v1.InventoryService/GetPart"}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))

	_, err := interceptor(ctx, nil, info, noopHandler(nil))

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestAuthInterceptor_success(t *testing.T) {
	iamClient := grpcMocks.NewMockIAMClient(t)
	interceptor := NewAuthInterceptor(iamClient)

	token := "session-token"
	user := model.User{Uuid: "user-uuid"}
	iamClient.EXPECT().Whoami(mock.Anything, token).Return(user, nil).Once()

	info := &grpc.UnaryServerInfo{FullMethod: "/inventory.v1.InventoryService/GetPart"}
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))

	var gotCtx context.Context
	resp, err := interceptor(ctx, nil, info, noopHandler(&gotCtx))

	assert.NoError(t, err)
	assert.Equal(t, "ok", resp)

	gotUser, ok := model.UserFromContext(gotCtx)
	assert.True(t, ok)
	assert.Equal(t, user, gotUser)
}
