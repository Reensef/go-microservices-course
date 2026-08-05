//go:build integration

package integration

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	iamV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

// testPassword — пароль, используемый для всех тестовых пользователей
const testPassword = "SuperSecret123!"

// newIamGRPCConn открывает соединение с контейнером iam, поднятым в тестовом окружении
func newIamGRPCConn() *grpc.ClientConn {
	conn, err := grpc.NewClient(
		env.IamApp.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	Expect(err).ToNot(HaveOccurred(), "expected successful connection to iam gRPC server")

	return conn
}

// loginRandomUser регистрирует случайного пользователя через iam.UserService.Register
// и логинит его через iam.AuthService.Login, возвращая session_uuid.
// Этот токен нужно передавать дальше в запросах к inventory через metadata "authorization".
func loginRandomUser(ctx context.Context) string {
	conn := newIamGRPCConn()

	userClient := iamV1.NewUserServiceClient(conn)
	authClient := iamV1.NewAuthServiceClient(conn)

	login := gofakeit.Username() + "_" + gofakeit.UUID()

	_, err := userClient.Register(ctx, &iamV1.RegisterRequest{
		Info: &iamV1.UserRegistrationInfo{
			Info: &iamV1.UserInfo{
				Login: login,
				Email: gofakeit.Email(),
			},
			Password: testPassword,
		},
	})
	Expect(err).ToNot(HaveOccurred(), "expected successful registration")

	loginResp, err := authClient.Login(ctx, &iamV1.LoginRequest{
		Login:    login,
		Password: testPassword,
	})
	Expect(err).ToNot(HaveOccurred(), "expected successful login")

	return loginResp.GetSessionUuid()
}

// withAuthorization прикладывает session_uuid как gRPC-metadata "authorization",
// имитируя то, как order прокидывает токен в inventory.
func withAuthorization(ctx context.Context, sessionUUID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+sessionUUID)
}
