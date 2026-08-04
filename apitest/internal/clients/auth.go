package clients

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	iamGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

const (
	bearerPrefix = "Bearer "
	testPassword = "SuperSecret123!"
)

type sessionCtxKey struct{}

// WithSession кладёт session_uuid в контекст сразу для обоих протоколов:
// в исходящую gRPC-metadata (читают inventory/payment) и в само значение
// контекста (читает bearerRoundTripper для HTTP-клиента order). Один и тот же
// ctx можно передавать в любой из клиентов Clients.
func WithSession(ctx context.Context, sessionUUID string) context.Context {
	ctx = context.WithValue(ctx, sessionCtxKey{}, sessionUUID)
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", bearerPrefix+sessionUUID)
	return ctx
}

func sessionFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(sessionCtxKey{}).(string)
	return token, ok
}

// RegisterAndLogin регистрирует нового случайного пользователя через
// iam.UserService и логинит его через iam.AuthService. Возвращает
// session_uuid — токен, которым дальше авторизуются все вызовы сценария.
func RegisterAndLogin(t *testing.T, ctx context.Context, c *Clients) (sessionUUID, login string) {
	t.Helper()

	login = "apitest-" + uuid.NewString()

	_, err := c.IAMUser.Register(ctx, &iamGrpc.RegisterRequest{
		Info: &iamGrpc.UserRegistrationInfo{
			Info: &iamGrpc.UserInfo{
				Login: login,
				Email: login + "@example.com",
			},
			Password: testPassword,
		},
	})
	require.NoError(t, err, "register user")

	loginResp, err := c.IAMAuth.Login(ctx, &iamGrpc.LoginRequest{
		Login:    login,
		Password: testPassword,
	})
	require.NoError(t, err, "login user")

	return loginResp.GetSessionUuid(), login
}
