package model

import "github.com/google/uuid"

type PartsFilter struct {
	Uuids                 []uuid.UUID
	Names                 []string
	Categories            []PartCategory
	ManufacturerCountries []string
	Tags                  []string
}
