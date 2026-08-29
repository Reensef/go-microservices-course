package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/order/internal/model"
)

func TestMiddleware_missingHeader(t *testing.T) {
	mw := NewAuthMiddleware()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mw(next).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.False(t, called)
}

func TestMiddleware_success(t *testing.T) {
	mw := NewAuthMiddleware()

	userID := "user-uuid"

	var gotUser model.User
	var ok bool

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-user-id", userID)
	rec := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser, ok = model.UserFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, ok)
	assert.Equal(t, model.User{Uuid: userID}, gotUser)
}
