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
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type PartManufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type PartInfo struct {
	Name          string                           `bson:"name"`
	Description   string                           `bson:"description"`
	Price         float64                          `bson:"price"`
	StockQuantity int64                            `bson:"stockQuantity"`
	Category      PartCategory                     `bson:"category"`
	Dimensions    *PartDimensions                  `bson:"dimensions"`
	Manufacturer  *PartManufacturer                `bson:"manufacturer"`
	Tags          []string                         `bson:"tags"`
	Metadata      map[string]multivalue.MultiValue `bson:"metadata"`
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
	Metadata      map[string]*multivalue.MultiValue
}

type Part struct {
	ID        string    `bson:"_id"`
	Info      PartInfo  `bson:"info"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}
