package part

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Reensef/go-microservices-course/inventory/internal/model"
	"github.com/Reensef/go-microservices-course/inventory/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
)

func TestGetByFilter_Uuid(t *testing.T) {
	repo := &repository{
		parts: make(map[uuid.UUID]*repoModel.Part, 10),
	}

	uuids := make([]uuid.UUID, 0, 10)
	filterUuids := make([]uuid.UUID, 0, 2)
	for range 2 {
		uuid := uuid.New()
		uuids = append(uuids, uuid)
		filterUuids = append(filterUuids, uuid)
	}

	for _, uuid := range uuids {
		repoModelOrder := &repoModel.Part{
			Uuid: uuid,
		}
		repo.parts[uuid] = repoModelOrder
	}

	filtered := repo.GetByFilter(t.Context(), &model.PartsFilter{
		Uuids: filterUuids,
	})

	for _, filteredOrder := range filtered {
		assert.True(t, filteredOrder.Uuid == filterUuids[0] ||
			filteredOrder.Uuid == filterUuids[1])
	}
}

func TestGetByFilter_Name(t *testing.T) {
	repo := &repository{
		parts: make(map[uuid.UUID]*repoModel.Part, 10),
	}

	names := []string{"name1", "name2", "name1", "name3", "name4"}
	filterNames := []string{"name1", "name2"}

	for _, name := range names {
		repoModelOrder := &repoModel.Part{
			Info: repoModel.PartInfo{
				Name: name,
			},
		}
		repo.parts[uuid.New()] = repoModelOrder
	}

	filtered := repo.GetByFilter(t.Context(), &model.PartsFilter{
		Names: filterNames,
	})

	for _, filteredOrder := range filtered {
		assert.True(t, filteredOrder.Info.Name == filterNames[0] ||
			filteredOrder.Info.Name == filterNames[1])
	}
}

func TestGetByFilter_Category(t *testing.T) {
	repo := &repository{
		parts: make(map[uuid.UUID]*repoModel.Part, 10),
	}

	categoryVariants := []repoModel.PartCategory{
		repoModel.PartCategory_ENGINE,
		repoModel.PartCategory_FUEL,
		repoModel.PartCategory_PORTHOLE,
		repoModel.PartCategory_WING,
	}
	categories := make([]repoModel.PartCategory, 0, 10)
	filterCategories := []model.PartCategory{}

	for range 2 {
		cat := categoryVariants[rand.IntN(len(categoryVariants))]
		categories = append(categories, cat)
		filterCategories = append(filterCategories, converter.ToModelPartCategory(cat))
	}

	for range 8 {
		cat := categoryVariants[rand.IntN(len(categoryVariants))]
		categories = append(categories, cat)
	}

	for _, category := range categories {
		repoModelOrder := &repoModel.Part{
			Info: repoModel.PartInfo{
				Category: category,
			},
		}
		repo.parts[uuid.New()] = repoModelOrder
	}

	filtered := repo.GetByFilter(t.Context(), &model.PartsFilter{
		Categories: filterCategories,
	})

	for _, filteredOrder := range filtered {
		assert.True(t, filteredOrder.Info.Category == filterCategories[0] ||
			filteredOrder.Info.Category == filterCategories[1])
	}
}

func TestGetByFilter_ManufacturerCountry(t *testing.T) {
	repo := &repository{
		parts: make(map[uuid.UUID]*repoModel.Part, 10),
	}

	variants := make([]string, 0, 4)
	for range cap(variants) {
		variants = append(variants, gofakeit.Company())
	}
	manufacturers := make([]string, 0, 10)
	filterManufacturers := make([]string, 0, 2)

	for range 2 {
		m := variants[rand.IntN(len(variants))]
		manufacturers = append(manufacturers, m)
		filterManufacturers = append(filterManufacturers, m)
	}

	for range 8 {
		m := variants[rand.IntN(len(variants))]
		manufacturers = append(manufacturers, m)
	}

	for _, man := range manufacturers {
		repoModelOrder := &repoModel.Part{
			Info: repoModel.PartInfo{
				Manufacturer: &repoModel.PartManufacturer{
					Country: man,
				},
			},
		}
		repo.parts[uuid.New()] = repoModelOrder
	}

	filtered := repo.GetByFilter(t.Context(), &model.PartsFilter{
		ManufacturerCountries: filterManufacturers,
	})

	for _, filteredOrder := range filtered {
		assert.True(t, filteredOrder.Info.Manufacturer.Country == filterManufacturers[0] ||
			filteredOrder.Info.Manufacturer.Country == filterManufacturers[1])
	}
}

func TestGetByFilter_Tags(t *testing.T) {
	repo := &repository{
		parts: make(map[uuid.UUID]*repoModel.Part, 10),
	}

	variants := make([]string, 0, 5)
	for range cap(variants) {
		variants = append(variants, gofakeit.Word())
	}
	data := make([]string, 0, 10)
	filterData := make([]string, 0, 2)

	for range 2 {
		m := variants[rand.IntN(len(variants))]
		data = append(data, m)
		filterData = append(filterData, m)
	}

	for range 8 {
		m := variants[rand.IntN(len(variants))]
		data = append(data, m)
	}

	for i := 1; i < len(data); i++ {
		repoModelOrder := &repoModel.Part{
			Info: repoModel.PartInfo{
				Tags: []string{data[i], data[i-1]},
			},
		}
		repo.parts[uuid.New()] = repoModelOrder
	}

	filtered := repo.GetByFilter(t.Context(), &model.PartsFilter{
		Tags: filterData,
	})

	for _, filteredOrder := range filtered {
		if !slices.Contains(filteredOrder.Info.Tags, filterData[0]) &&
			!slices.Contains(filteredOrder.Info.Tags, filterData[1]) {
			assert.Fail(t, "tag not found")
		}
	}
}
