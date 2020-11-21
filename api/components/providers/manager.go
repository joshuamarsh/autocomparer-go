package providers

import (
	"carcompare/structs"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/apsdehal/go-logger"
)

// Manager handles data from multiple providers
type Manager struct {
	providers map[string]Provider
	logger    *logger.Logger
}

// NewManager helper function for creating a new manager
func NewManager(l *logger.Logger) (Manager, error) {
	m := Manager{
		providers: make(map[string]Provider),
	}

	m.logger = l

	return m, nil
}

// RegisterProvider store a new provider in the manager
func (m *Manager) RegisterProvider(id string, p Provider) error {
	m.logger.Infof("Registering new provider %s", id)
	if _, exists := m.providers[id]; exists {
		return fmt.Errorf("Provider %s already exists", id)
	}

	m.providers[id] = p
	return nil
}

// GetAdvert sorts all adverts from all sources into array
func (m *Manager) GetAdvert(providers []string, brand string, model string, postcode string, radius string, sortBy string) (structs.AdvertProviders, error) {
	for _, provider := range providers {
		if _, ok := m.providers[provider]; !ok {
			return structs.AdvertProviders{}, fmt.Errorf("Provider %s is not supported", provider)
		}
	}

	providerLock := sync.Mutex{}
	providerResponses := make(chan map[string]structs.Adverts, len(providers))
	providerErrors := make(chan error, len(providers))

	var advert structs.AdvertProviders

	for _, provider := range providers {
		go func(provider string) {
			res, err := m.providers[provider].GetAdvert(postcode, radius, brand, model, sortBy)
			if err != nil {
				providerErrors <- err
			} else {
				responses := make(map[string]structs.Adverts)
				for _, r := range res {
					responses[provider] = r
				}

				providerResponses <- responses
			}
		}(provider)
	}

	advertProviders := []structs.Advert{}
	for _, provider := range providers {
		m.logger.Debugf(provider + " success")
		select {
		case err := <-providerErrors:
			return structs.AdvertProviders{}, err
		case res := <-providerResponses:
			providerLock.Lock()
			for _, r := range res {
				advertProviders = append(advertProviders, r.Adverts...)
			}
			providerLock.Unlock()
		}
	}

	switch sortBy {
	case "best_match":
		sortBy = "BestMatch"
	case "date_desc":
		sortBy = "StartTimeNewest"
	case "date_asc":
		sortBy = ""
	case "dist_asc":
		sortBy = "DistanceNearest"
	case "year_desc":
		sortBy = "year-desc"
	case "price_asc":
		sort.Slice(advertProviders, func(i, j int) bool {
			return advertProviders[i].Price < advertProviders[j].Price
		})
	case "price_desc":
		sort.Slice(advertProviders, func(i, j int) bool {
			return advertProviders[i].Price > advertProviders[j].Price
		})
	case "miles_asc":
		sortBy = ""
	default:
		sortBy = ""
	}

	advert.Adverts = advertProviders
	advert.Providers = providers

	return advert, nil
}

func (m *Manager) GetMakes(providers []string) ([]structs.MakeProvider, error) {
	for _, provider := range providers {
		if _, ok := m.providers[provider]; !ok {
			return []structs.MakeProvider{}, fmt.Errorf("Provider %s is not supported", provider)
		}
	}

	providerLock := sync.Mutex{}
	providerResponses := make(chan map[string][]structs.Make, len(providers))
	providerErrors := make(chan error, len(providers))

	for _, provider := range providers {
		go func(provider string) {
			res, err := m.providers[provider].GetMakes()
			if err != nil {
				providerErrors <- err
			} else {
				responses := make(map[string][]structs.Make)
				responses[provider] = res

				providerResponses <- responses
			}
		}(provider)
	}

	categoryProviders := []structs.MakeProvider{}
	for _, provider := range providers {
		m.logger.Debugf(provider + " success")
		select {
		case err := <-providerErrors:
			return []structs.MakeProvider{}, err
		case res := <-providerResponses:
			providerLock.Lock()
			for _, r := range res {
				categories := []structs.MakeProvider{}
				for _, category := range r {
					newCategory := true
					for index, categoryProvider := range categoryProviders {
						if category.Name == categoryProvider.Name {
							categoryProviders[index].Providers[category.Provider] = category.ID
							newCategory = false
						}
					}
					if newCategory {
						id := strings.ToLower(strings.ReplaceAll(category.Name, " ", ""))
						categories = append(categories, structs.MakeProvider{ID: id, Providers: map[string]string{category.Provider: category.ID}, Name: category.Name})
					}
				}
				categoryProviders = append(categoryProviders, categories...)
			}
			providerLock.Unlock()
		}
	}

	return categoryProviders, nil
}

func (m *Manager) GetModels(providers []string, brand string) ([]structs.ModelProvider, error) {
	for _, provider := range providers {
		if _, ok := m.providers[provider]; !ok {
			return []structs.ModelProvider{}, fmt.Errorf("Provider %s is not supported", provider)
		}
	}

	providerLock := sync.Mutex{}
	providerResponses := make(chan map[string][]structs.Model, len(providers))
	providerErrors := make(chan error, len(providers))

	for _, provider := range providers {
		go func(provider string) {
			res, err := m.providers[provider].GetModels(brand)
			if err != nil {
				providerErrors <- err
			} else {
				responses := make(map[string][]structs.Model)
				responses[provider] = res

				providerResponses <- responses
			}
		}(provider)
	}

	modelProviders := []structs.ModelProvider{}
	for _, provider := range providers {
		m.logger.Debugf(provider + " success")
		select {
		case err := <-providerErrors:
			return []structs.ModelProvider{}, err
		case res := <-providerResponses:
			providerLock.Lock()
			for _, r := range res {
				models := []structs.ModelProvider{}
				for _, model := range r {
					newModel := true
					for index, modelProvider := range modelProviders {
						if model.Name == modelProvider.Name {
							modelProviders[index].Providers[model.Provider] = model.ID
							newModel = false
						}
					}
					if newModel {
						id := strings.ToLower(strings.ReplaceAll(model.Name, " ", ""))
						models = append(models, structs.ModelProvider{ID: id, Providers: map[string]string{model.Provider: model.ID}, Name: model.Name})
					}
				}
				modelProviders = append(modelProviders, models...)
			}
			providerLock.Unlock()
		}
	}

	return modelProviders, nil
}
