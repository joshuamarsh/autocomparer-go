package providers

import "carcompare/structs"

// Provider handles data from a single source
type Provider interface {
	GetAdvert(postcode string, radius string, make string, model string, sortBy string) ([]structs.Adverts, error)
	GetMakes() ([]structs.Make, error)
	GetModels(brand string) ([]structs.Model, error)
}
