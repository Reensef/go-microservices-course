//go:build integration

package integration

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	iamV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

// testPassword — пароль, используемый для всех тестовых пользователей
const testPassword = "SuperSecret123!"

// newGRPCConn открывает соединение с приложением, поднятым в тестовом окружении
func newGRPCConn() *grpc.ClientConn {
	conn, err := grpc.NewClient(
		env.App.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	Expect(err).ToNot(HaveOccurred(), "expected successful connection to gRPC server")

	return conn
}

// registerRandomUser регистрирует случайного пользователя через UserService.Register
// и возвращает его uuid, логин и пароль — общие тестовые данные для Auth- и User-тестов.
func registerRandomUser(ctx context.Context, userClient iamV1.UserServiceClient) (userUUID, login, password string) {
	login = gofakeit.Username() + "_" + gofakeit.UUID()
	password = testPassword

	resp, err := userClient.Register(ctx, &iamV1.RegisterRequest{
		Info: &iamV1.UserRegistrationInfo{
			Info: &iamV1.UserInfo{
				Login: login,
				Email: gofakeit.Email(),
			},
			Password: password,
		},
	})
	Expect(err).ToNot(HaveOccurred(), "expected successful registration")

	return resp.GetUserUuid(), login, password
}
