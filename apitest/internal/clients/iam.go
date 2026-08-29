package clients

import (
	"context"
	"net/http"

	iamGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/iam/v1"
)

// iamClient — HTTP-клиент iam через Envoy-гейтвей (POST /iam/v1/users,
// POST /iam/v1/sessions), а не нативный gRPC — так сценарии проходят через
// тот же путь (ext_authz, grpc_json_transcoder), что и реальные клиенты.
type iamClient struct {
	baseURL string
	http    *http.Client
}

func (c *iamClient) Register(ctx context.Context, req *iamGrpc.RegisterRequest) (*iamGrpc.RegisterResponse, error) {
	resp := &iamGrpc.RegisterResponse{}
	if err := doJSON(ctx, c.http, http.MethodPost, c.baseURL+"/iam/v1/users", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *iamClient) Login(ctx context.Context, req *iamGrpc.LoginRequest) (*iamGrpc.LoginResponse, error) {
	resp := &iamGrpc.LoginResponse{}
	if err := doJSON(ctx, c.http, http.MethodPost, c.baseURL+"/iam/v1/sessions", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
