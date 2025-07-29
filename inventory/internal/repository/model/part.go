package model

import (
	"time"

	"github.com/google/uuid"
)

type PartCategory int32

const (
	PartCategory_UNSPECIFIED PartCategory = iota // Неизвестная категория
	PartCategory_ENGINE                          // Двигатель
	PartCategory_FUEL                            // Топливо
	PartCategory_PORTHOLE                        // Иллюминатор
	PartCategory_WING                            // Крыло
)

type PartDimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type PartManufacturer struct {
	Name    string
	Country string
	Website string
}

type PartInfo struct {
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      PartCategory
	Dimensions    *PartDimensions
	Manufacturer  *PartManufacturer
	Tags          []string
	Metadata      map[string]*any
}

type PartUpdateInfo struct {
	Name          *string
	Description   *string
	Price         *float64
	StockQuantity *int64
	Category      *PartCategory
	Dimensions    *PartDimensions
	Manufacturer  *PartManufacturer
	Tags          []string
	Metadata      map[string]*any
}

type Part struct {
	Uuid      uuid.UUID
	Info      PartInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}
