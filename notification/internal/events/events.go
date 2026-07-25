package events

import "context"

type Consumer interface {
	RunConsumer(ctx context.Context) error
}
