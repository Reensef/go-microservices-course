package v1

import (
	"context"
	"errors"
	"log"
	"strings"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"

	"github.com/Reensef/go-microservices-course/iam/internal/model"
)

const bearerPrefix = "Bearer "

// Check реализует envoy.service.auth.v3.Authorization/Check: вызывается фильтром
// ext_authz Envoy на каждый запрос к защищённому роуту. Проверка сессии выполняется
// in-process — напрямую через service.AuthService.Whoami, без сетевого похода,
// так как этот обработчик живёт внутри самого IAM.
func (a *api) Check(
	ctx context.Context,
	req *authv3.CheckRequest,
) (*authv3.CheckResponse, error) {
	headers := req.GetAttributes().GetRequest().GetHttp().GetHeaders()

	authHeader, ok := headers["authorization"]
	if !ok || !strings.HasPrefix(authHeader, bearerPrefix) {
		return deniedResponse(typev3.StatusCode_Unauthorized, "missing or malformed authorization header"), nil
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token == "" {
		return deniedResponse(typev3.StatusCode_Unauthorized, "missing or malformed authorization header"), nil
	}

	_, user, err := a.service.Whoami(ctx, token)
	if err != nil {
		log.Printf("authz api: error getting whoami: %s", err.Error())

		switch {
		case errors.Is(err, model.ErrSessionUuidInvalidFormat),
			errors.Is(err, model.ErrSessionNotFound),
			errors.Is(err, model.ErrUserNotFound):
			return deniedResponse(typev3.StatusCode_Unauthorized, "invalid or expired session"), nil
		default:
			return deniedResponse(typev3.StatusCode_Forbidden, "internal server error"), nil
		}
	}

	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(code.Code_OK)},
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers: []*corev3.HeaderValueOption{
					{
						Header: &corev3.HeaderValue{
							Key:   "x-user-id",
							Value: user.Uuid,
						},
					},
				},
			},
		},
	}, nil
}

func deniedResponse(httpStatus typev3.StatusCode, body string) *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(code.Code_PERMISSION_DENIED)},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{Code: httpStatus},
				Body:   body,
			},
		},
	}
}
