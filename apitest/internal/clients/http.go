package clients

import "net/http"

// bearerRoundTripper прикладывает Authorization-заголовок к каждому HTTP-запросу
// ogen-клиента order, читая токен из контекста запроса — тот же контекст, что
// используется для проброса токена в gRPC-metadata через WithSession, см. auth.go.
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
