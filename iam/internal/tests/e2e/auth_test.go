//go:build integration

package integration

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	iamV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

var _ = Describe("AuthService", func() {
	var (
		ctx        context.Context
		cancel     context.CancelFunc
		authClient iamV1.AuthServiceClient
		userClient iamV1.UserServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		conn := newGRPCConn()
		authClient = iamV1.NewAuthServiceClient(conn)
		userClient = iamV1.NewUserServiceClient(conn)
	})

	AfterEach(func() {
		err := env.ClearState(ctx)
		Expect(err).ToNot(HaveOccurred(), "expected to successfully clear test state")

		cancel()
	})

	Describe("Login", func() {
		It("must return a session uuid for a registered user with valid credentials", func() {
			_, login, password := registerRandomUser(ctx, userClient)

			resp, err := authClient.Login(ctx, &iamV1.LoginRequest{
				Login:    login,
				Password: password,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(uuid.Validate(resp.GetSessionUuid())).To(Succeed())
		})

		It("must return Unauthenticated for a wrong password", func() {
			_, login, _ := registerRandomUser(ctx, userClient)

			_, err := authClient.Login(ctx, &iamV1.LoginRequest{
				Login:    login,
				Password: "definitely-wrong-password",
			})

			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.Unauthenticated))
		})

		It("must return Unauthenticated for a non-existent login", func() {
			_, err := authClient.Login(ctx, &iamV1.LoginRequest{
				Login:    "no-such-user-" + uuid.NewString(),
				Password: testPassword,
			})

			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.Unauthenticated))
		})
	})

	Describe("Whoami", func() {
		It("must return the session and its owner for an active session", func() {
			userUUID, login, password := registerRandomUser(ctx, userClient)

			loginResp, err := authClient.Login(ctx, &iamV1.LoginRequest{
				Login:    login,
				Password: password,
			})
			Expect(err).ToNot(HaveOccurred())

			resp, err := authClient.Whoami(ctx, &iamV1.WhoamiRequest{
				SessionUuid: loginResp.GetSessionUuid(),
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetSession().GetUuid()).To(Equal(loginResp.GetSessionUuid()))
			Expect(resp.GetUser().GetUuid()).To(Equal(userUUID))
			Expect(resp.GetUser().GetInfo().GetLogin()).To(Equal(login))
		})

		It("must return InvalidArgument for a malformed session uuid", func() {
			_, err := authClient.Whoami(ctx, &iamV1.WhoamiRequest{
				SessionUuid: "not-a-uuid",
			})

			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.InvalidArgument))
		})

		It("must return NotFound for a well-formed but unknown session uuid", func() {
			_, err := authClient.Whoami(ctx, &iamV1.WhoamiRequest{
				SessionUuid: uuid.NewString(),
			})

			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.NotFound))
		})

		It("must return NotFound when the session's owner no longer exists", func() {
			userUUID, login, password := registerRandomUser(ctx, userClient)

			loginResp, err := authClient.Login(ctx, &iamV1.LoginRequest{
				Login:    login,
				Password: password,
			})
			Expect(err).ToNot(HaveOccurred())

			err = env.DeleteUser(ctx, userUUID)
			Expect(err).ToNot(HaveOccurred(), "expected to successfully delete the user directly from postgres")

			_, err = authClient.Whoami(ctx, &iamV1.WhoamiRequest{
				SessionUuid: loginResp.GetSessionUuid(),
			})

			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.NotFound))
		})
	})
})
