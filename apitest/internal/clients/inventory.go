package clients

import (
	"context"
	"net/http"
	"net/url"

	inventoryGrpc "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

// inventoryClient — HTTP-клиент inventory через Envoy-гейтвей
// (POST /inventory/v1/parts:search, GET /inventory/v1/parts/{id}).
type inventoryClient struct {
	baseURL string
	http    *http.Client
}

func (c *inventoryClient) ListParts(ctx context.Context, req *inventoryGrpc.ListPartsRequest) (*inventoryGrpc.ListPartsResponse, error) {
	resp := &inventoryGrpc.ListPartsResponse{}
	if err := doJSON(ctx, c.http, http.MethodPost, c.baseURL+"/inventory/v1/parts:search", req.GetFilter(), resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *inventoryClient) GetPart(ctx context.Context, req *inventoryGrpc.GetPartRequest) (*inventoryGrpc.GetPartResponse, error) {
	resp := &inventoryGrpc.GetPartResponse{}
	target := c.baseURL + "/inventory/v1/parts/" + url.PathEscape(req.GetId())
	if err := doJSON(ctx, c.http, http.MethodGet, target, nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
