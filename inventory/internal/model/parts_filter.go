package model

type PartsFilter struct {
	IDs                   []string
	Names                 []string
	Categories            []PartCategory
	ManufacturerCountries []string
	Tags                  []string
}
