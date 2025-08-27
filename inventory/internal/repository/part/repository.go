package part

import (
	"context"
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/mongo"

	model "github.com/Reensef/go-microservices-course/inventory/internal/model"
)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	collection := db.Collection("parts")
	repo := &repository{
		collection: collection,
	}

	repo.genData()

	return repo
}

func (s *repository) genData() {
	for range 10 {
		part, err := s.Create(context.Background(), generateRandomPart())
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(part.ID)
	}
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
