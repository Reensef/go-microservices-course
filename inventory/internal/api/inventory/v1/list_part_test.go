package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	"github.com/Reensef/go-microservices-course/inventory/internal/service/mocks"
	inventoryV1 "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

func TestListPart(t *testing.T) {
	service := mocks.NewMockInventoryService(t)

	dataLen := 5

	service.EXPECT().GetPartsByFilter(t.Context(), mock.Anything).Return(make([]*model.Part, dataLen), nil).Once()

	a := NewAPI(service)

	response, err := a.ListParts(t.Context(), &inventoryV1.ListPartsRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, dataLen, len(response.Parts))
}
