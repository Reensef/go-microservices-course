package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Reensef/go-microservices-course/iam/internal/model"
	iamV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

func ToProtoUser(user model.User) *iamV1.User {
	return &iamV1.User{
		Uuid:      user.Uuid,
		Info:      ToProtoUserInfo(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func ToProtoUserInfo(info model.UserInfo) *iamV1.UserInfo {
	return &iamV1.UserInfo{
		Login:               info.Login,
		Email:               info.Email,
		NotificationMethods: ToProtoNotificationMethods(info.NotificationMethods),
	}
}

func ToProtoNotificationMethods(methods []model.NotificationMethod) []*iamV1.NotificationMethod {
	result := make([]*iamV1.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		result = append(result, &iamV1.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}
	return result
}

// info передаётся как proto-указатель, а не значение: генерируемые protobuf-геттеры
// (GetLogin, GetEmail, ...) nil-safe, а сама структура iamV1.UserInfo в Go всегда передаётся как *T.
func ToModelUserInfo(info *iamV1.UserInfo) model.UserInfo {
	return model.UserInfo{
		Login:               info.GetLogin(),
		Email:               info.GetEmail(),
		NotificationMethods: ToModelNotificationMethods(info.GetNotificationMethods()),
	}
}

func ToModelNotificationMethods(methods []*iamV1.NotificationMethod) []model.NotificationMethod {
	result := make([]model.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		result = append(result, model.NotificationMethod{
			ProviderName: method.GetProviderName(),
			Target:       method.GetTarget(),
		})
	}
	return result
}
