package service

import (
	"context"
	"log"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	repo "github.com/Reensef/go-microservices-course/inventory/internal/repository"
)

type service struct {
	repo repo.PartRepository
}

func NewService(repo repo.PartRepository) *service {
	service := &service{
		repo: repo,
	}

	return service
}

func (s *service) GenData() error {
	for range 10 {
		part, err := s.repo.Create(context.Background(), generateRandomPart())
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(part.ID)
	}
	return nil
}

func generateRandomPart() *model.PartInfo {
	// Генерация случайной категории
	categories := []model.PartCategory{
		model.PartCategory_UNSPECIFIED,
		model.PartCategory_ENGINE,
		model.PartCategory_FUEL,
		model.PartCategory_PORTHOLE,
		model.PartCategory_WING,
	}
	category := categories[gofakeit.Number(0, len(categories)-1)]

	// Генерация случайных тегов
	tags := gofakeit.RandomString([]string{"metal", "plastic", "electronics", "rare", "expensive"})

	partInfo := &model.PartInfo{
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Sentence(10),
		Price:         float64(gofakeit.Price(100, 10000)),
		StockQuantity: int64(gofakeit.Number(0, 1000)),
		Category:      category,
		Dimensions: &model.PartDimensions{
			Length: gofakeit.Float64Range(1, 100),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: gofakeit.Float64Range(0.1, 500),
		},
		Manufacturer: &model.PartManufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: []string{tags},
	}

	return partInfo
}
