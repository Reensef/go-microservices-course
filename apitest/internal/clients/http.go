package clients

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// bearerRoundTripper прикладывает Authorization-заголовок к каждому HTTP-запросу,
// читая токен из контекста запроса — см. WithSession в auth.go. Используется
// всеми HTTP-клиентами (iam, inventory, order), все они бьют в один и тот же
// Envoy-гейтвей.
type bearerRoundTripper struct {
	base http.RoundTripper
}

func (rt bearerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	base := rt.base
	if base == nil {
		base = http.DefaultTransport
	}

	token, ok := sessionFromContext(req.Context())
	if !ok || token == "" {
		return base.RoundTrip(req)
	}

	// http.RoundTripper не должен модифицировать исходный запрос.
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", bearerPrefix+token)

	return base.RoundTrip(req)
}

// doJSON шлёт reqMsg как JSON-тело (protojson, с сохранением proto-имён полей —
// так же, как настроен grpc_json_transcoder в envoy.yaml) и разбирает JSON-ответ
// в respMsg. reqMsg/respMsg могут быть nil (запрос без тела / ответ без тела).
func doJSON(ctx context.Context, httpClient *http.Client, method, url string, reqMsg, respMsg proto.Message) error {
	var body io.Reader
	if reqMsg != nil {
		b, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(reqMsg)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		body = bytes.NewReader(b)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	if body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = httpResp.Body.Close() }() //nolint:gosec,errcheck

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if httpResp.StatusCode >= 300 {
		return fmt.Errorf("%s %s: status %d: %s", method, url, httpResp.StatusCode, string(respBytes))
	}

	if respMsg != nil && len(respBytes) > 0 {
		if err := protojson.Unmarshal(respBytes, respMsg); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}
