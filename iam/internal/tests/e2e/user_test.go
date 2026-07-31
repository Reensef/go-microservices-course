//go:build integration

package integration

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	iamV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

var _ = Describe("UserService", func() {
	var (
		ctx        context.Context
		cancel     context.CancelFunc
		userClient iamV1.UserServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		conn := newGRPCConn()
		userClient = iamV1.NewUserServiceClient(conn)
	})

	AfterEach(func() {
		err := env.ClearState(ctx)
		Expect(err).ToNot(HaveOccurred(), "expected to successfully clear test state")

		cancel()
	})

	Describe("Register", func() {
		It("must register a new user and return a valid uuid", func() {
			login := gofakeit.Username() + "_" + uuid.NewString()
			email := gofakeit.Email()

			resp, err := userClient.Register(ctx, &iamV1.RegisterRequest{
				Info: &iamV1.UserRegistrationInfo{
					Info: &iamV1.UserInfo{
						Login: login,
						Email: email,
					},
					Password: testPassword,
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(uuid.Validate(resp.GetUserUuid())).To(Succeed())
		})

		It("must return AlreadyExists when the login is taken", func() {
			_, login, _ := registerRandomUser(ctx, userClient)

			_, err := userClient.Register(ctx, &iamV1.RegisterRequest{
				Info: &iamV1.UserRegistrationInfo{
					Info: &iamV1.UserInfo{
						Login: login,
						Email: gofakeit.Email(),
					},
					Password: testPassword,
				},
			})

			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.AlreadyExists))
		})
	})

	Describe("GetUser", func() {
		It("must return an existing user by uuid", func() {
			userUUID, login, _ := registerRandomUser(ctx, userClient)

			resp, err := userClient.GetUser(ctx, &iamV1.GetUserRequest{
				UserUuid: userUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetUser().GetUuid()).To(Equal(userUUID))
			Expect(resp.GetUser().GetInfo().GetLogin()).To(Equal(login))
		})

		It("must return InvalidArgument for a malformed uuid", func() {
			_, err := userClient.GetUser(ctx, &iamV1.GetUserRequest{
				UserUuid: "not-a-uuid",
			})

			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})

		It("must return NotFound for a well-formed but unknown uuid", func() {
			_, err := userClient.GetUser(ctx, &iamV1.GetUserRequest{
				UserUuid: uuid.NewString(),
			})

			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.NotFound))
		})
	})
})
