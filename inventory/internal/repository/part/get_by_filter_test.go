package part

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	converter "github.com/Reensef/go-microservices-course/inventory/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
)

func TestGetByFilter_Uuid(t *testing.T) {
	repo := &repository{
		parts: make(map[uuid.UUID]*repoModel.Part),
	}

	repoModelOrder := generateRandomPart()
	modelOrder := converter.RepoModelPartToModel(repoModelOrder)

	fakeModelOrder := converter.RepoModelPartToModel(generateRandomPart())

	repo.parts[repoModelOrder.Uuid] = repoModelOrder

	actual, err := repo.GetByUuid(t.Context(), &repoModelOrder.Uuid)

	assert.Equal(t, modelOrder, actual)
	assert.NotEqual(t, fakeModelOrder, actual)
	assert.Nil(t, err)

	fakeUuid := uuid.New()
	actual, err = repo.GetByUuid(t.Context(), &fakeUuid)
	assert.Nil(t, actual)
	assert.ErrorIs(t, model.ErrPartNotFound, err)
}
