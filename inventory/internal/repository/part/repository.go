package part

import (
	"log"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
)

type repository struct {
	mu    sync.RWMutex
	parts map[uuid.UUID]*repoModel.Part
}

func NewRepository() *repository {
	repo := &repository{
		parts: make(map[uuid.UUID]*repoModel.Part),
	}

	for range 10 {
		fake := generateRandomPart()
		log.Println(fake.Uuid.String())
		repo.parts[fake.Uuid] = fake
	}

	return repo
}

func generateRandomPart() *repoModel.Part {
	// Генерация случайной категории
	categories := []repoModel.PartCategory{
		repoModel.PartCategory_UNSPECIFIED,
		repoModel.PartCategory_ENGINE,
		repoModel.PartCategory_FUEL,
		repoModel.PartCategory_PORTHOLE,
		repoModel.PartCategory_WING,
	}
	category := categories[gofakeit.Number(0, len(categories)-1)]

	// Генерация случайных тегов
	tags := gofakeit.RandomString([]string{"metal", "plastic", "electronics", "rare", "expensive"})

	// Создание детали
	part := &repoModel.Part{
		Uuid: uuid.MustParse(gofakeit.UUID()),
		Info: repoModel.PartInfo{
			Name:          gofakeit.ProductName(),
			Description:   gofakeit.Sentence(10),
			Price:         float64(gofakeit.Price(100, 10000)),
			StockQuantity: int64(gofakeit.Number(0, 1000)),
			Category:      category,
			Dimensions: &repoModel.PartDimensions{
				Length: gofakeit.Float64Range(1, 100),
				Width:  gofakeit.Float64Range(1, 100),
				Height: gofakeit.Float64Range(1, 100),
				Weight: gofakeit.Float64Range(0.1, 500),
			},
			Manufacturer: &repoModel.PartManufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags: []string{tags},
		},
		CreatedAt: gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
		UpdatedAt: time.Now(),
	}

	return part
}
