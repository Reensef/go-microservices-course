package converter

import (
	model "github.com/Reensef/go-microservices-course/iam/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/iam/internal/repository/model"
)

func ToModelUser(user repoModel.User) model.User {
	return model.User{
		Uuid:      user.Uuid,
		Info:      ToModelUserInfo(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToModelUserInfo(info repoModel.UserInfo) model.UserInfo {
	return model.UserInfo{
		Login:               info.Login,
		Email:               info.Email,
		NotificationMethods: ToModelNotificationMethods(info.NotificationMethods),
	}
}

func ToModelNotificationMethods(methods []repoModel.NotificationMethod) []model.NotificationMethod {
	result := make([]model.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		result = append(result, model.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}
	return result
}

func ToRepoNotificationMethods(methods []model.NotificationMethod) []repoModel.NotificationMethod {
	result := make([]repoModel.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		result = append(result, repoModel.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}
	return result
}
