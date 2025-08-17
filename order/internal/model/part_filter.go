package model

type PartsFilter struct {
	Ids                   []string
	Names                 []string
	Categories            []PartCategory
	ManufacturerCountries []string
	Tags                  []string
}
