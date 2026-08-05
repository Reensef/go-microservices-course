package middleware

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	grpcClients "github.com/Reensef/go-microservices-course/order/internal/client/grpc"
	"github.com/Reensef/go-microservices-course/order/internal/model"
)

const bearerPrefix = "Bearer "

type errorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(errorBody{Code: status, Message: message}); err != nil {
		log.Printf("middleware: failed to write error response: %s", err)
	}
}

func NewAuthMiddleware(iamClient grpcClients.IAMClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, bearerPrefix) {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			token := strings.TrimPrefix(header, bearerPrefix)
			if token == "" {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			user, err := iamClient.Whoami(r.Context(), token)
			if err != nil {
				if errors.Is(err, model.ErrInvalidSession) {
					writeError(w, http.StatusUnauthorized, "unauthorized")
					return
				}

				writeError(w, http.StatusInternalServerError, "internal server error")
				return
			}

			ctx := model.WithToken(model.WithUser(r.Context(), user), token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
