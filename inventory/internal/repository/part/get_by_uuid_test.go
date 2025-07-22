package part

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
)

// TestGetByUuid_ExistsUuid проверяет случай, когда деталь существует в репозитории,
// и ее uuid можно найти.
func TestGetByUuid_ExistsUuid(t *testing.T) {
	repo := &repository{
		parts: make(map[uuid.UUID]*repoModel.Part),
	}

	uuid := uuid.New()

	repo.parts[uuid] = &repoModel.Part{}

	actual, err := repo.GetByUuid(t.Context(), uuid)

	assert.NotNil(t, actual)
	assert.NoError(t, err)
}
