package service

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	"github.com/Reensef/go-microservices-course/inventory/internal/repository/mocks"
)

// TestGetByUuid_Exists проверяет случай, когда деталь существует в репозитории.
func TestGetByUuid_Exists(t *testing.T) {
	repoMock := mocks.NewMockPartRepository(t)

	uuid := uuid.New()

	repoMock.EXPECT().GetByUuid(t.Context(), uuid).Return(&model.Part{}, nil).Once()

	service := NewService(repoMock)

	part, err := service.GetPartByUuid(t.Context(), uuid)

	assert.NoError(t, err)
	assert.NotNil(t, part)
}

// TestGetByUuid_Exists проверяет случай, когда деталь отсутствует в репозитории.
func TestGetByUuid_NotExists(t *testing.T) {
	repoMock := mocks.NewMockPartRepository(t)

	uuid := uuid.New()

	repoMock.EXPECT().GetByUuid(t.Context(), uuid).Return(nil, fmt.Errorf("error")).Once()

	service := NewService(repoMock)

	part, err := service.GetPartByUuid(t.Context(), uuid)

	assert.Error(t, err)
	assert.Nil(t, part)
}
