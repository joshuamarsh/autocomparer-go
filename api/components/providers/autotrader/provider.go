package autotrader

import (
	"bytes"
	"carcompare/api/database"
	"carcompare/api/models"
	"carcompare/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/apsdehal/go-logger"
	"github.com/gocolly/colly"
)

// Provider handles communications with a single service
type Provider struct {
	logger *logger.Logger
}

// NewProvider helper function for creating a advert source
func NewProvider(l *logger.Logger) *Provider {
	return &Provider{
		logger: l,
	}
}

// GetAdvert gets adverts from autotrader
func (p *Provider) GetAdvert(postcode string, radius string, brand string, model string, sortBy string) ([]structs.Adverts, error) {
	p.logger.Notice("GetAdvert for Autotrader")

	autotraderSort := ""
	switch sortBy {
	case "best_match":
		autotraderSort = "relevance"
	case "date_desc":
		autotraderSort = "datedesc"
	case "date_asc":
		autotraderSort = ""
	case "dist_asc":
		autotraderSort = "distance"
	case "year_desc":
		autotraderSort = "year-desc"
	case "price_asc":
		autotraderSort = "price-asc"
	case "price_desc":
		autotraderSort = "price-desc"
	case "miles_asc":
		autotraderSort = "mileage"
	default:
		autotraderSort = ""
	}

	baseURL := "https://www.autotrader.co.uk"

	db := database.DB
	if brand != "" {
		var makeAutotrader models.Make
		db.Where("value = ? AND provider = ?", brand, "autotrader").First(&makeAutotrader)
		if makeAutotrader.ProviderValue != "" {
			brand = url.QueryEscape(makeAutotrader.ProviderValue)
		}
	}

	if model != "" {
		var modelAutotrader models.Model
		db.Where("value = ? AND provider = ? AND make = ?", model, "autotrader", brand).First(&modelAutotrader)
		if modelAutotrader.ProviderValue != "" {
			model = url.QueryEscape(modelAutotrader.ProviderValue)
		}
	}

	urlQuery := fmt.Sprintf("/car-search?sort=%s&postcode=%s&radius=%s&make=%s&model=%s", autotraderSort, postcode, radius, brand, model)
	p.logger.Debugf("%s", urlQuery)
	c := colly.NewCollector()

	Adverts := structs.Adverts{}
	c.OnHTML(".search-page__result", func(e *colly.HTMLElement) {
		featuredListing := e.Attr("data-is-featured-listing")
		promotedListing := e.Attr("data-is-promoted-listing")
		ymalListing := e.Attr("data-is-ymal-listing")

		if featuredListing != "true" && promotedListing != "true" && ymalListing != "true" {
			price := e.ChildText(".physical-stock-now")
			if price == "" {
				price = e.ChildText(".vehicle-price")
			}
			if price == "" {
				price = e.ChildText(".product-card-pricing__price span")
			}

			priceReplace := strings.NewReplacer("£", "", ",", "", ".", "")
			priceFloat := priceReplace.Replace(price)
			priceFormatted, _ := strconv.ParseInt(priceFloat, 10, 64)

			title := e.ChildText(".product-card-details__title")
			if title == "" {
				title = e.ChildText(".listing-title a")
			}

			image := e.ChildAttr(".product-card-image__main-image", "src")
			if image == "" {
				image = e.ChildAttr(".listing-main-image img", "src")
			}

			linkParams := e.ChildAttr(".listing-fpa-link", "href")
			if linkParams == "" {
				linkParams = "/car-details/" + e.Attr("id")
			}

			mileage := ""
			e.ForEachWithBreak(".listing-key-specs li", func(_ int, specElement *colly.HTMLElement) bool {
				if strings.Contains(specElement.Text, "miles") {
					mileage = specElement.Text
					return false
				}
				return true
			})

			description := e.ChildText(".listing-description")
			if description == "" {
				description = e.ChildText(".product-card-details__attention-grabber")
			}

			distanceString := strings.Replace(strings.Replace(e.Attr("data-distance-value"), "s", "", -1), " mile", "", 1)
			distance, _ := strconv.ParseUint(strings.Replace(distanceString, " miles", "", 1), 0, 64)

			if e.Attr("id") != "" && title != "" && price != "" {
				advert := structs.Advert{
					Provider:    "autotrader",
					ID:          e.Attr("id"),
					Link:        "https://www.autotrader.co.uk" + linkParams,
					Location:    e.ChildText(".seller-town"),
					Distance:    distance,
					Title:       title,
					Price:       priceFormatted * 100,
					Mileage:     mileage,
					Description: description,
					Image:       image,
				}
				Adverts.AddAdvert(advert)
			}
		}
	})

	c.Visit(baseURL + urlQuery)
	// c.Visit(baseURL + urlQuery + "&page=2")
	// c.Visit(baseURL + urlQuery + "&page=3")
	// c.Visit(baseURL + urlQuery + "&page=4")
	// c.Visit(baseURL + urlQuery + "&page=5")
	// fmt.Printf("%v", Adverts)

	providerAdverts := make([]structs.Adverts, 0)
	providerAdverts = append(providerAdverts, Adverts)
	return providerAdverts, nil
}

func (p *Provider) GetMakes() ([]structs.Make, error) {
	p.logger.Notice("GetMakes for Autotrader")

	baseURL := "https://www.autotrader.co.uk"

	c := colly.NewCollector()

	categories := []structs.Make{}
	// db := database.DB
	c.OnHTML(".atc-field__input.atc-field__input--select optgroup[label='All makes'] option", func(e *colly.HTMLElement) {
		name := strings.Split(e.Text, " (")[0]
		switch name {
		case "Vauxhall":
			name = "Vauxhall/Opel"
		case "MINI":
			name = "Mini"
		case "KIA":
			name = "Kia"
		case "Land Rover":
			name = "Land Rover/Range Rover"
		case "SKODA":
			name = "Skoda"
		case "Custom Vehicle":
			name = "Kit Cars/Custom Vehicle"
		}

		value := e.Attr("value")
		// makeDB := models.Make{
		// 	Value:    value,
		// 	Name:     name,
		// 	Provider: "autotrader",
		// }
		// if err := db.Create(&makeDB).Error; err != nil {
		// 	p.logger.Errorf("failed to create make %s", name)
		// }

		category := structs.Make{
			Provider: "autotrader",
			ID:       value,
			Name:     name,
		}

		categories = append(categories, category)
	})

	c.Visit(baseURL)

	return categories, nil
}

func (p *Provider) GetModels(brand string) ([]structs.Model, error) {
	p.logger.Notice("GetModels for Autotrader")

	db := database.DB
	var make models.Make
	db.Where("value = ? AND provider = ?", brand, "autotrader").First(&make)

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	getModelsRequest := GetModelsRequest{
		OperationName: "SearchFormFacetsQuery",
		Variables: Variables{
			AdvertQuery: AdvertQuery{
				Make: []string{make.ProviderValue},
			},
			Facets: []string{"model"},
		},
		Query: "query SearchFormFacetsQuery($advertQuery: AdvertQuery!, $facets: [SearchFacetName]) { search { adverts(advertQuery: $advertQuery) { facets(facets: $facets) { name values { name value } } } }}",
	}

	modelsRequestBody, err := json.Marshal(getModelsRequest)

	req, err := http.NewRequest("POST", "https://www.autotrader.co.uk/at-graphql", bytes.NewBuffer(modelsRequestBody))
	if err != nil {
		return []structs.Model{}, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return []structs.Model{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []structs.Model{}, err
	}

	var modelsResponse GetModelsResponse
	if err := json.Unmarshal(body, &modelsResponse); err != nil {
		return []structs.Model{}, err
	}

	if len(modelsResponse.Data.Search.Adverts.Facets) < 1 {
		return []structs.Model{}, nil
	}

	providerModels := []structs.Model{}
	// db := database.DB
	for _, modelResponse := range modelsResponse.Data.Search.Adverts.Facets[0].Values {
		modelName := strings.Title(strings.ToLower(modelResponse.Name))
		modelValue := modelResponse.Value
		if modelName != "" && modelValue != "" {
			model := structs.Model{
				Provider: "autotrader",
				ID:       modelValue,
				Name:     modelName,
			}
			providerModels = append(providerModels, model)
		}
	}

	return providerModels, nil
}
