package clients

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	iamGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

const (
	bearerPrefix = "Bearer "
	testPassword = "SuperSecret123!"
)

type sessionCtxKey struct{}

// WithSession кладёт session_uuid в контекст — bearerRoundTripper читает его
// оттуда и подставляет в заголовок Authorization для всех HTTP-клиентов.
func WithSession(ctx context.Context, sessionUUID string) context.Context {
	return context.WithValue(ctx, sessionCtxKey{}, sessionUUID)
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

	_, err := c.IAM.Register(ctx, &iamGrpc.RegisterRequest{
		Info: &iamGrpc.UserRegistrationInfo{
			Info: &iamGrpc.UserInfo{
				Login: login,
				Email: login + "@example.com",
			},
			Password: testPassword,
		},
	})
	require.NoError(t, err, "register user")

	loginResp, err := c.IAM.Login(ctx, &iamGrpc.LoginRequest{
		Login:    login,
		Password: testPassword,
	})
	require.NoError(t, err, "login user")

	return loginResp.GetSessionUuid(), login
}
