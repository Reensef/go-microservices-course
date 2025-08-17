package v1

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	"github.com/Reensef/go-microservices-course/inventory/internal/service/mocks"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func TestGetPart_ValidUuid(t *testing.T) {
	service := mocks.NewMockInventoryService(t)

	uuid := uuid.NewString()

	service.EXPECT().GetPartByID(t.Context(), uuid).Return(&model.Part{}, nil).Once()

	a := NewAPI(service)

	response, err := a.GetPart(t.Context(), &inventoryV1.GetPartRequest{Id: uuid})
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.Part)
}

func TestGetPart_ServiceError(t *testing.T) {
	service := mocks.NewMockInventoryService(t)

	uuid := uuid.NewString()

	service.EXPECT().GetPartByID(t.Context(), uuid).Return(nil, fmt.Errorf("error")).Once()

	a := NewAPI(service)

	response, err := a.GetPart(t.Context(), &inventoryV1.GetPartRequest{Id: uuid})
	assert.Error(t, err)
	assert.Nil(t, response)
}
