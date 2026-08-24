package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	grpcMocks "github.com/Reensef/go-microservices-course/order/internal/client/grpc/mocks"
	"github.com/Reensef/go-microservices-course/order/internal/model"
)

func TestMiddleware_missingHeader(t *testing.T) {
	iamClient := grpcMocks.NewMockIAMClient(t)
	mw := NewAuthMiddleware(iamClient)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mw(next).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.False(t, called)
	assert.Empty(t, iamClient.Calls)
}

func TestMiddleware_malformedHeader(t *testing.T) {
	iamClient := grpcMocks.NewMockIAMClient(t)
	mw := NewAuthMiddleware(iamClient)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Token abc")
	rec := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next should not be called")
	})).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Empty(t, iamClient.Calls)
}

func TestMiddleware_invalidSession(t *testing.T) {
	iamClient := grpcMocks.NewMockIAMClient(t)
	mw := NewAuthMiddleware(iamClient)

	token := "session-token"
	iamClient.EXPECT().Whoami(mock.Anything, token).Return(model.User{}, model.ErrInvalidSession).Once()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next should not be called")
	})).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestMiddleware_iamError(t *testing.T) {
	iamClient := grpcMocks.NewMockIAMClient(t)
	mw := NewAuthMiddleware(iamClient)

	token := "session-token"
	iamClient.EXPECT().Whoami(mock.Anything, token).Return(model.User{}, fmt.Errorf("unavailable")).Once()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next should not be called")
	})).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestMiddleware_success(t *testing.T) {
	iamClient := grpcMocks.NewMockIAMClient(t)
	mw := NewAuthMiddleware(iamClient)

	token := "session-token"
	user := model.User{Uuid: "user-uuid"}
	iamClient.EXPECT().Whoami(mock.Anything, token).Return(user, nil).Once()

	var gotUser model.User
	var ok bool

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser, ok = model.UserFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, ok)
	assert.Equal(t, user, gotUser)
}
