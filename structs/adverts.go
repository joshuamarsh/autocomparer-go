package structs

import (
	"crypto/sha256"
	"fmt"
)

// AdvertProviders stores adverts from all providers
type AdvertProviders struct {
	Providers                  []string `json:"providers"`
	LowestProviderHighestPrice *int64   `json:"lowest_provider_highest_price"`
	HighestProviderLowestPrice *int64   `json:"highest_provider_lowest_price"`
	Adverts                    []Advert `json:"adverts"`
}

// Adverts stores adverts from a single provider
type Adverts struct {
	Adverts                    []Advert `json:"adverts"`
	LowestProviderHighestPrice *int64
	HighestProviderLowestPrice *int64
}

// Advert store information about a single advert
type Advert struct {
	Provider    string `json:"provider"`
	ID          string `json:"id"`
	Link        string `json:"link"`
	Location    string `json:"location"`
	Distance    uint64 `json:"distance"`
	Title       string `json:"title"`
	Price       int64  `json:"price"`
	Mileage     string `json:"mileage"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

// Location stores location parameters
type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// Make stores information about a make from a provider
type Make struct {
	Provider string `json:"provider"`
	ID       string `json:"id"`
	Name     string `json:"name"`
}

// MakeProvider stores information about a make
type MakeProvider struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Model stores information about a model from a provider
type Model struct {
	Provider string `json:"provider"`
	ID       string `json:"id"`
	Name     string `json:"name"`
}

// ModelProvider stores information about a model
type ModelProvider struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Sha256 generates unique token
func Sha256(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
