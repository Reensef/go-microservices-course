package middleware

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	"github.com/Reensef/go-microservices-course/platform/pkg/logger"
)

type errorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(errorBody{Code: status, Message: message}); err != nil {
		logger.Error("middleware: failed to write error response", zap.Error(err))
	}
}

func NewAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Header.Get("x-user-id")
			if userID == "" {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			ctx := model.WithUser(r.Context(), model.User{Uuid: userID})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
