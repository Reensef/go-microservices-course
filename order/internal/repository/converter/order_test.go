package converter

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/order/internal/model"
	repoModel "github.com/Reensef/go-microservices-course/order/internal/repository/model"
)

func TestToRepoModelOrderInfo(t *testing.T) {
	tests := []struct {
		name     string
		info     *model.OrderInfo
		expected *repoModel.OrderInfo
	}{
		{
			name: "Valid OrderInfo",
			info: &model.OrderInfo{
				UserUuid: "ea815eb8-a569-47fb-89fb-3052b3612396",

				PartIds: []string{
					"ea815eb8-a569-47fb-89fb-3052b3612396",
					"2f217a02-22c4-4044-b5a3-7eb06001a521",
					"3f74bbbe-ad48-43cd-bffd-80d0a1b6a8e6",
				},
				TransactionUuid: "2f217a02-22c4-4044-b5a3-7eb06001a521",
				TotalPrice:      10.99,
				PaymentMethod:   model.OrderPaymentMethod_CARD,
				Status:          model.OrderStatus_PENDING_PAYMENT,
			},
			expected: &repoModel.OrderInfo{
				UserUuid: "ea815eb8-a569-47fb-89fb-3052b3612396",
				PartUuids: []string{
					"ea815eb8-a569-47fb-89fb-3052b3612396",
					"2f217a02-22c4-4044-b5a3-7eb06001a521",
					"3f74bbbe-ad48-43cd-bffd-80d0a1b6a8e6",
				},
				TransactionUuid: sql.NullString{
					String: "2f217a02-22c4-4044-b5a3-7eb06001a521",
					Valid:  true,
				},
				TotalPrice:    10.99,
				PaymentMethod: repoModel.OrderPaymentMethod_CARD,
				Status:        repoModel.OrderStatus_PENDING_PAYMENT,
			},
		},
		{
			name:     "Nil OrderInfo",
			info:     nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ToRepoModelOrderInfo(tt.info)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestToRepoModelOrder(t *testing.T) {
	tests := []struct {
		name     string
		order    *model.Order
		expected *repoModel.Order
	}{
		{
			name: "Valid Order",
			order: &model.Order{
				Uuid:      "ea815eb8-a569-47fb-89fb-3052b3612396",
				Info:      model.OrderInfo{},
				CreatedAt: time.Date(2021, 1, 1, 1, 1, 1, 1, time.UTC),
				UpdatedAt: time.Date(2022, 2, 2, 1, 2, 2, 2, time.UTC),
				DeletedAt: func() *time.Time {
					t := time.Date(2023, 3, 3, 3, 3, 3, 3, time.UTC)

					return &t
				}(),
			},
			expected: &repoModel.Order{
				Uuid:      "ea815eb8-a569-47fb-89fb-3052b3612396",
				Info:      repoModel.OrderInfo{},
				CreatedAt: time.Date(2021, 1, 1, 1, 1, 1, 1, time.UTC),
				UpdatedAt: time.Date(2022, 2, 2, 1, 2, 2, 2, time.UTC),
				DeletedAt: func() *time.Time {
					t := time.Date(2023, 3, 3, 3, 3, 3, 3, time.UTC)

					return &t
				}(),
			},
		},
		{
			name:     "Nil Order",
			order:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ToRepoModelOrder(tt.order)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestToModelOrderInfo(t *testing.T) {
	tests := []struct {
		name     string
		info     *repoModel.OrderInfo
		expected *model.OrderInfo
	}{
		{
			name: "Valid OrderInfo",
			info: &repoModel.OrderInfo{
				UserUuid: "ea815eb8-a569-47fb-89fb-3052b3612396",
				PartUuids: []string{
					"ea815eb8-a569-47fb-89fb-3052b3612396",
					"2f217a02-22c4-4044-b5a3-7eb06001a521",
					"3f74bbbe-ad48-43cd-bffd-80d0a1b6a8e6",
				},
				TransactionUuid: sql.NullString{
					String: "2f217a02-22c4-4044-b5a3-7eb06001a521",
					Valid:  true,
				},
				TotalPrice:    10.99,
				PaymentMethod: repoModel.OrderPaymentMethod_CARD,
				Status:        repoModel.OrderStatus_PENDING_PAYMENT,
			},
			expected: &model.OrderInfo{
				UserUuid: "ea815eb8-a569-47fb-89fb-3052b3612396",
				PartIds: []string{
					"ea815eb8-a569-47fb-89fb-3052b3612396",
					"2f217a02-22c4-4044-b5a3-7eb06001a521",
					"3f74bbbe-ad48-43cd-bffd-80d0a1b6a8e6",
				},
				TransactionUuid: "2f217a02-22c4-4044-b5a3-7eb06001a521",
				TotalPrice:      10.99,
				PaymentMethod:   model.OrderPaymentMethod_CARD,
				Status:          model.OrderStatus_PENDING_PAYMENT,
			},
		},

		{
			name:     "Nil OrderInfo",
			info:     nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ToModelOrderInfo(tt.info)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestToModelOrder(t *testing.T) {
	tests := []struct {
		name     string
		order    *repoModel.Order
		expected *model.Order
	}{
		{
			name: "Valid Order",
			order: &repoModel.Order{
				Uuid:      "ea815eb8-a569-47fb-89fb-3052b3612396",
				Info:      repoModel.OrderInfo{},
				CreatedAt: time.Date(2021, 1, 1, 1, 1, 1, 1, time.UTC),
				UpdatedAt: time.Date(2022, 2, 2, 1, 2, 2, 2, time.UTC),
				DeletedAt: func() *time.Time {
					t := time.Date(2023, 3, 3, 3, 3, 3, 3, time.UTC)

					return &t
				}(),
			},

			expected: &model.Order{
				Uuid:      "ea815eb8-a569-47fb-89fb-3052b3612396",
				Info:      model.OrderInfo{},
				CreatedAt: time.Date(2021, 1, 1, 1, 1, 1, 1, time.UTC),
				UpdatedAt: time.Date(2022, 2, 2, 1, 2, 2, 2, time.UTC),
				DeletedAt: func() *time.Time {
					t := time.Date(2023, 3, 3, 3, 3, 3, 3, time.UTC)

					return &t
				}(),
			},
		},

		{
			name:     "Nil Order",
			order:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ToModelOrder(tt.order)

			assert.Equal(t, tt.expected, actual)
		})
	}
}
