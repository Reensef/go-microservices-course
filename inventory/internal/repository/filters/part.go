package filters

import (
	"slices"

	model "github.com/Reensef/go-microservices-course/inventory/internal/model"
	"github.com/Reensef/go-microservices-course/inventory/internal/repository/converter"
	repoModel "github.com/Reensef/go-microservices-course/inventory/internal/repository/model"
)

// Проверяет, соответствует ли деталь всем условиям фильтра
func MatchPartFilters(part *repoModel.Part, filter *model.PartsFilter) bool {
	if !filterPartByNames(part, filter.Names) {
		return false
	}
	if !filterPartByCategories(part, filter.Categories) {
		return false
	}
	if !filterPartByManufacturerCountries(part, filter.ManufacturerCountries) {
		return false
	}
	if !filterPartByTags(part, filter.Tags) {
		return false
	}
	return true
}

// Проверяет, соответствует ли имя детали списку имен из фильтра
func filterPartByNames(part *repoModel.Part, names []string) bool {
	if len(names) == 0 {
		return true
	}
	return slices.Contains(names, part.Info.Name)
}

// Проверяет, соответствует ли категория детали списку категорий из фильтра
func filterPartByCategories(part *repoModel.Part, categories []model.PartCategory) bool {
	if len(categories) == 0 {
		return true
	}
	return slices.Contains(
		categories,
		converter.RepoModelPartCategoryToModel(part.Info.Category),
	)
}

// Проверяет, соответствует ли страна производителя списку стран из фильтра
func filterPartByManufacturerCountries(part *repoModel.Part, countries []string) bool {
	if len(countries) == 0 {
		return true
	}
	return slices.Contains(countries, part.Info.Manufacturer.Country)
}

// Проверяет, содержит ли деталь хотя бы один тег из списка тегов в фильтре
func filterPartByTags(part *repoModel.Part, tags []string) bool {
	if len(tags) == 0 {
		return true
	}
	for _, tag := range tags {
		if slices.Contains(part.Info.Tags, tag) {
			return true
		}
	}
	return false
}
