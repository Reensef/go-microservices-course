package notification

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	clientMocks "github.com/Reensef/go-microservices-course/notification/internal/client/telegram/mocks"
	"github.com/Reensef/go-microservices-course/notification/internal/model"
)

func TestNotifyOrderPaid_success(t *testing.T) {
	telegramClient := clientMocks.NewMockClient(t)
	chatID := int64(12345)
	s := New(telegramClient, chatID)

	event := model.OrderPaidEvent{
		UUID:            uuid.NewString(),
		OrderUUID:       uuid.NewString(),
		UserUUID:        uuid.NewString(),
		TransactionUUID: uuid.NewString(),
		PaymentMethod:   "CARD",
	}

	telegramClient.EXPECT().
		SendMessage(context.Background(), chatID, mock.Anything).
		Return(nil).
		Once()

	err := s.NotifyOrderPaid(context.Background(), event)

	assert.NoError(t, err)
}

func TestNotifyOrderPaid_clientError(t *testing.T) {
	telegramClient := clientMocks.NewMockClient(t)
	chatID := int64(12345)
	s := New(telegramClient, chatID)

	event := model.OrderPaidEvent{
		OrderUUID: uuid.NewString(),
	}
	sendError := fmt.Errorf("telegram error")

	telegramClient.EXPECT().
		SendMessage(context.Background(), chatID, mock.Anything).
		Return(sendError).
		Once()

	err := s.NotifyOrderPaid(context.Background(), event)

	assert.Equal(t, sendError, err)
}
