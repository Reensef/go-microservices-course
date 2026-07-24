package service

import (
	"context"
)

type AssemblyService interface {
	AssembleShip(ctx context.Context, orderUuid, userUuid string) error
}
