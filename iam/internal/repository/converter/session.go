package converter

import (
	model "github.com/Reensef/go-microservices-course/iam/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/iam/internal/repository/model"
)

func ToModelSession(session repoModel.Session) model.Session {
	return model.Session{
		Uuid:      session.Uuid,
		UserUuid:  session.UserUuid,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
		ExpiresAt: session.ExpiresAt,
	}
}
