package main

import (
	"slices"

	inventory "github.com/Reensef/go-microservices-course/shared/pkg/proto/inventory/v1"
)

// Проверяет, соответствует ли деталь всем условиям фильтра
func matchPartFilters(part *inventory.Part, filter *inventory.PartsFilter) bool {
	if !filterPartByNames(part, filter.GetNames()) {
		return false
	}
	if !filterPartByCategories(part, filter.GetCategories()) {
		return false
	}
	if !filterPartByManufacturerCountries(part, filter.GetManufacturerCountries()) {
		return false
	}
	if !filterPartByTags(part, filter.GetTags()) {
		return false
	}
	return true
}

// Проверяет, соответствует ли имя детали списку имен из фильтра
func filterPartByNames(part *inventory.Part, names []string) bool {
	if len(names) == 0 {
		return true
	}
	return slices.Contains(names, part.Name)
}

// Проверяет, соответствует ли категория детали списку категорий из фильтра
func filterPartByCategories(part *inventory.Part, categories []inventory.Category) bool {
	if len(categories) == 0 {
		return true
	}
	return slices.Contains(categories, part.Category)
}

// Проверяет, соответствует ли страна производителя списку стран из фильтра
func filterPartByManufacturerCountries(part *inventory.Part, countries []string) bool {
	if len(countries) == 0 {
		return true
	}
	return slices.Contains(countries, part.Manufacturer.Country)
}

// Проверяет, содержит ли деталь хотя бы один тег из списка тегов в фильтре
func filterPartByTags(part *inventory.Part, tags []string) bool {
	if len(tags) == 0 {
		return true
	}
	for _, tag := range tags {
		if slices.Contains(part.Tags, tag) {
			return true
		}
	}
	return false
}
