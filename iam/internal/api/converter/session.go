package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Reensef/go-microservices-course/iam/internal/model"
	iamV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

func ToProtoSession(session model.Session) *iamV1.Session {
	return &iamV1.Session{
		Uuid:      session.Uuid,
		CreatedAt: timestamppb.New(session.CreatedAt),
		UpdatedAt: timestamppb.New(session.UpdatedAt),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}
