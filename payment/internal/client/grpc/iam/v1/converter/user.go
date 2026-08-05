package converter

import (
	"github.com/Reensef/go-microservices-course/payment/internal/model"
	iamGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

func ToModelUser(user *iamGrpc.User) model.User {
	return model.User{
		Uuid:  user.GetUuid(),
		Login: user.GetInfo().GetLogin(),
		Email: user.GetInfo().GetEmail(),
	}
}
