package order

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func TestCreateOrder(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := &repository{
			data: make(map[uuid.UUID]*repoModel.Order),
		}

		uuid, err := repo.CreateOrder(context.Background(), &model.OrderInfo{})
		require.NoError(t, err)
		require.NotNil(t, uuid)
		assert.Len(t, repo.data, 1)
	})

	t.Run("Nil input", func(t *testing.T) {
		repo := &repository{
			data: make(map[uuid.UUID]*repoModel.Order),
		}
		uuid, err := repo.CreateOrder(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, uuid)
		assert.Len(t, repo.data, 0)
	})

	t.Run("Concurrent creation", func(t *testing.T) {
		repo := &repository{
			data: make(map[uuid.UUID]*repoModel.Order),
		}
		info := &model.OrderInfo{}

		var wg sync.WaitGroup
		for range 10 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := repo.CreateOrder(context.Background(), info)
				require.NoError(t, err)
			}()
		}
		wg.Wait()
		assert.Len(t, repo.data, 10)
	})
}
