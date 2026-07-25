package kafka

import (
	"slices"

	"github.com/IBM/sarama"
)

// Middleware — функция middleware для дополнительной обработки.
type Middleware func(next MessageHandler) MessageHandler

// groupHandler — обёртка для sarama.ConsumerGroupHandler
type groupHandler struct {
	handler MessageHandler
}

// NewGroupHandler создаёт новый groupHandler с middleware цепочкой.
func NewGroupHandler(handler MessageHandler, middlewares ...Middleware) *groupHandler {
	// Применяем middleware цепочку
	for _, middleware := range slices.Backward(middlewares) {
		handler = middleware(handler)
	}

	return &groupHandler{
		handler: handler,
	}
}

func (g *groupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g *groupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			msg := Message{
				Key:            message.Key,
				Value:          message.Value,
				Topic:          message.Topic,
				Partition:      message.Partition,
				Offset:         message.Offset,
				Timestamp:      message.Timestamp,
				BlockTimestamp: message.BlockTimestamp,
				Headers:        extractHeaders(message.Headers),
			}

			if err := g.handler(session.Context(), msg); err != nil {
				continue
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func extractHeaders(headers []*sarama.RecordHeader) map[string][]byte {
	result := make(map[string][]byte)
	for _, h := range headers {
		if h != nil && h.Key != nil {
			result[string(h.Key)] = h.Value
		}
	}

	return result
}
