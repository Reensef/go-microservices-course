package model

import (
	"time"

	"github.com/Reensef/go-microservices-course/platform/pkg/multivalue"
)

type PartCategory int32

const (
	PartCategory_UNSPECIFIED PartCategory = 0 // Неизвестная категория
	PartCategory_ENGINE      PartCategory = 1 // Двигатель
	PartCategory_FUEL        PartCategory = 2 // Топливо
	PartCategory_PORTHOLE    PartCategory = 3 // Иллюминатор
	PartCategory_WING        PartCategory = 4 // Крыло
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
	Metadata      map[string]multivalue.MultiValue
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
	Id        string
	Info      PartInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}
