package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	"github.com/Reensef/go-microservices-course/inventory/internal/repository/mocks"
)

// TestGetByFilter проверяет, что метод GetPartsByFilter сервиса
// корректно возвращает ожидаемое количество деталей при взаимодействии с
// фиктивным репозиторием.
func TestGetByFilter(t *testing.T) {
	repoMock := mocks.NewMockPartRepository(t)

	dataLen := 10

	repoMock.EXPECT().GetByFilter(t.Context(), mock.Anything).Return(make([]*model.Part, dataLen), nil).Once()

	service := NewService(repoMock)

	parts, err := service.GetPartsByFilter(t.Context(), nil)

	assert.NoError(t, err)
	assert.Equal(t, dataLen, len(parts))
}
