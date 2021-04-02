package facebook

import (
	"bytes"
	"carcompare/api/database"
	"carcompare/api/models"
	"carcompare/config"
	"carcompare/structs"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/apsdehal/go-logger"
	"gorm.io/gorm"
)

// Provider handles communications with a single service
type Provider struct {
	logger *logger.Logger
}

// NewProvider helper function for creating a new advert source
func NewProvider(l *logger.Logger) *Provider {
	return &Provider{
		logger: l,
	}
}

// GetAdvert gets adverts from facebook
func (p *Provider) GetAdvert(postcode string, radius string, brand string, model string, sortBy string, page *uint) ([]structs.Advert, error) {
	p.logger.Notice("GetAdvert for Facebook")

	// facebookSort := ""
	// facebookSortOrder := ""
	// switch sortBy {
	// case "best_match":
	// 	facebookSort = ""
	// 	facebookSortOrder = ""
	// case "date_asc":
	// 	facebookSort = "CREATION_TIME"
	// 	facebookSortOrder = "ASCEND"
	// case "date_desc":
	// 	facebookSort = "CREATION_TIME"
	// 	facebookSortOrder = "DESCEND"
	// case "dist_asc":
	// 	facebookSort = "DISTANCE"
	// 	facebookSortOrder = "ASCEND"
	// case "year_asc":
	// 	facebookSort = "VEHICLE_YEAR"
	// 	facebookSortOrder = "ASCEND"
	// case "year_desc":
	// 	facebookSort = "VEHICLE_YEAR"
	// 	facebookSortOrder = "DESCEND"
	// case "price_asc":
	// 	facebookSort = "PRICE_AMOUNT"
	// 	facebookSortOrder = "ASCEND"
	// case "price_desc":
	// 	facebookSort = "PRICE_AMOUNT"
	// 	facebookSortOrder = "DESCEND"
	// case "miles_asc":
	// 	facebookSort = "VEHICLE_MILEAGE"
	// 	facebookSortOrder = "ASCEND"
	// case "miles_desc":
	// 	facebookSort = "VEHICLE_MILEAGE"
	// 	facebookSortOrder = "DESCEND"
	// default:
	// 	facebookSort = ""
	// 	facebookSortOrder = ""
	// }

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	db := database.DB
	searchBrand := brand
	if brand != "" {
		var makeFacebook models.Make
		db.Where("value = ? AND provider = ?", brand, "facebook").First(&makeFacebook)
		if makeFacebook.ProviderValue != "" {
			brand = makeFacebook.ProviderValue
		}
	}

	if model != "" {
		var modelFacebook models.Model
		db.Where("value = ? AND provider = ? AND make = ?", model, "facebook", searchBrand).First(&modelFacebook)
		if modelFacebook.ProviderValue != "" {
			model = modelFacebook.ProviderValue
		}
	}

	facebookAdverts := []Edge{}
	if page == nil || *page == 0 || *page == 1 {
		// params.Add("paginationInput.pageNumber", "1")
		getModelsRequest := GetAdvertsRequest{
			BuyLocation: BuyLocation{
				Latitude:  53.478698730469,
				Longitude: -2.2467041015625,
			},
			CategoryIDArray: []int64{807311116002614},
			ContextualData: []ContextualData{
				{Name: "make", Value: brand},
				{Name: "model", Value: model},
			},
			Count:                    24,
			MarketplaceBrowseContext: "CATEGORY_FEED",
			PriceRange:               []int64{0, 214748364700},
			Radius:                   60000,
			Scale:                    2,
			StringVerticalFields: []ContextualData{
				{Name: "canonical_make_id", Value: brand},
				{Name: "canonical_model_id", Value: model},
			},
			TopicPageParams: TopicPageParams{
				LocationID: "manchester",
				URL:        "vehicles",
			},
		}

		// if postcode != "" {
		// 	location, err := getLocation(postcode)
		// 	if err != nil {
		// 		return []structs.Adverts{}, err
		// 	}

		// 	getModelsRequest.BuyLocation = BuyLocation{
		// 		Latitude:  location.Latitude,
		// 		Longitude: location.Longitude,
		// 	}
		// }
		// getModelsRequest.BuyLocation = BuyLocation{
		// 	Latitude:  53.478698730469,
		// 	Longitude: -2.2467041015625,
		// }

		// p.logger.Debugf("%v", getModelsRequest.BuyLocation)

		// if brand != "" {
		// 	getModelsRequest.StringVerticalFields = append(getModelsRequest.StringVerticalFields, StringVerticalFields{Name: "canonical_make_id", Value: brand})
		// }

		// if brand != "" && model != "" {
		// 	getModelsRequest.StringVerticalFields = append(getModelsRequest.StringVerticalFields, StringVerticalFields{Name: "canonical_model_id", Value: model})
		// }

		// if facebookSort != "" {
		// 	getModelsRequest.FilterSortingParams = &FilterSortingParams{
		// 		SortByFilter: facebookSort,
		// 		SortOrder:    facebookSortOrder,
		// 	}
		// }

		modelsRequestBody, err := json.Marshal(getModelsRequest)
		if err != nil {
			return []structs.Advert{}, err
		}

		formData := url.Values{
			"fb_api_caller_class":      {"RelayModern"},
			"fb_api_req_friendly_name": {"CometMarketplaceCategoryContentContainerQuery"},
			"variables":                {string(modelsRequestBody)},
			"doc_id":                   {"5111225608947754"},
		}

		// p.logger.Debugf("%s", string(modelsRequestBody))

		req, err := http.NewRequest("POST", "https://en-gb.facebook.com/api/graphql/", strings.NewReader(formData.Encode()))
		if err != nil {
			return []structs.Advert{}, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		res, err := client.Do(req)
		if err != nil {
			return []structs.Advert{}, err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []structs.Advert{}, err
		}

		// p.logger.Debugf("%s", string(body))

		var advertsResponse GetAdvertsResponse
		if err := json.Unmarshal(body, &advertsResponse); err != nil {
			return []structs.Advert{}, err
		}

		if len(advertsResponse.Data.Viewer.MarketplaceFeedStories.Edges) < 1 {
			return []structs.Advert{}, nil
		}

		facebookAdverts = append(facebookAdverts, advertsResponse.Data.Viewer.MarketplaceFeedStories.Edges...)

		getAdvertsPaginationQueryParams := GetAdvertsPaginationQueryParams{
			postcode:  postcode,
			radius:    radius,
			brand:     brand,
			model:     model,
			sortBy:    sortBy,
			ItemIndex: 24,
			Count:     6,
		}
		facebookAdvertsFromQuery, err := getAdvertsPaginationQuery(getAdvertsPaginationQueryParams)
		if err != nil {
			return []structs.Advert{}, err
		}

		facebookAdverts = append(facebookAdverts, facebookAdvertsFromQuery...)
	} else {
		// ebayPage := strconv.Itoa(int(*page))
		// if ebayPage == "" {
		// 	params.Add("paginationInput.pageNumber", "1")
		// }
		// params.Add("paginationInput.pageNumber", ebayPage)
		itemIndex := (int64(*page) - 1) * 30

		getAdvertsPaginationQueryParams := GetAdvertsPaginationQueryParams{
			postcode:  postcode,
			radius:    radius,
			brand:     brand,
			model:     model,
			sortBy:    sortBy,
			ItemIndex: itemIndex,
			Count:     24,
		}
		facebookAdvertsFromQuery, err := getAdvertsPaginationQuery(getAdvertsPaginationQueryParams)
		if err != nil {
			return []structs.Advert{}, err
		}

		facebookAdverts = append(facebookAdverts, facebookAdvertsFromQuery...)

		getAdvertsPaginationQueryParams2 := GetAdvertsPaginationQueryParams{
			postcode:  postcode,
			radius:    radius,
			brand:     brand,
			model:     model,
			sortBy:    sortBy,
			ItemIndex: itemIndex + 24,
			Count:     6,
		}
		facebookAdvertsFromQuery2, err := getAdvertsPaginationQuery(getAdvertsPaginationQueryParams2)
		if err != nil {
			return []structs.Advert{}, err
		}

		facebookAdverts = append(facebookAdverts, facebookAdvertsFromQuery2...)
	}
	adverts := []structs.Advert{}
	for _, advertResponse := range facebookAdverts {
		priceReplace := strings.NewReplacer("£", "", ",", "", ".", "")
		priceFloat := priceReplace.Replace(advertResponse.Node.Listing.ListingPrice.FormattedAmount)
		priceFormatted, _ := strconv.ParseInt(priceFloat, 10, 64)
		// p.logger.DebugF("Ranking: %s", advertResponse.Node.Tracking)

		// var distance uint64
		// cityLocation, err := getCityLocation(advertResponse.Node.Listing.Location.ReverseGeocode.CityPage.DisplayName)
		// if err != nil {
		// 	p.logger.Errorf("%v", err)
		// } else {
		// 	distance = haversine(getModelsRequest.BuyLocation.Longitude, getModelsRequest.BuyLocation.Latitude, cityLocation.Longitude, cityLocation.Latitude)
		// }

		mileage := strings.Replace(strings.Replace(advertResponse.Node.Listing.CustomSubTitlesWithRenderingFlags[0].Subtitle, "K", ",000", 1), " · Dealership", "", 1)
		advert := structs.Advert{
			Provider:    "facebook",
			ID:          advertResponse.Node.StoryKey,
			Link:        "https://en-gb.facebook.com/marketplace/item/" + advertResponse.Node.Listing.ID,
			Location:    advertResponse.Node.Listing.Location.ReverseGeocode.City,
			Distance:    0,
			Title:       advertResponse.Node.Listing.MarketplaceListingTitle,
			Price:       priceFormatted * 100,
			Mileage:     mileage,
			Description: "",
			Image:       advertResponse.Node.Listing.PrimaryListingPhoto.Image.URI,
		}

		adverts = append(adverts, advert)
	}

	return adverts, nil
}

// GetMakes gets makes from facebook
func (p *Provider) GetMakes() ([]structs.Make, error) {
	p.logger.Notice("GetMakes for Facebook")

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	getMakesRequest := GetMakesRequest{
		CategoryIDS:            []int64{807311116002614},
		CategoryRankingEnabled: true,
		ContextualData: []ContextualData{
			{Name: "make", Value: "2260278494029417"},
			{Name: "seo_url", Value: "vehicles"},
		},
		HideL2Cats: true,
		Surface:    "CATEGORY_FEED",
		TopicPageParams: TopicPageParams{
			LocationID: "manchester",
			URL:        "vehicles",
		},
	}

	makesRequestBody, err := json.Marshal(getMakesRequest)
	if err != nil {
		return []structs.Make{}, err
	}

	formData := url.Values{
		"fb_api_caller_class":      {"RelayModern"},
		"fb_api_req_friendly_name": {"CometMarketplaceSearchRootQuery"},
		"variables":                {string(makesRequestBody)},
		"doc_id":                   {"5183985588339655"},
	}

	// p.logger.Debugf("%s", string(makesRequestBody))

	req, err := http.NewRequest("POST", "https://en-gb.facebook.com/api/graphql/", strings.NewReader(formData.Encode()))
	if err != nil {
		return []structs.Make{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return []structs.Make{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []structs.Make{}, err
	}

	// p.logger.Debugf("%s", string(body))

	var makesResponse GetMakesResponse
	if err := json.Unmarshal(body, &makesResponse); err != nil {
		return []structs.Make{}, err
	}

	makes := []structs.Make{}
	for _, filterType := range makesResponse.Data.Viewer.MarketplaceStructuredFields {
		if filterType.FilterKey == "make" {
			for _, makeResponse := range filterType.Choices {
				name := makeResponse.DisplayName
				switch name {
				case "Škoda":
					name = "Skoda"
				}
				make := structs.Make{
					Provider: "facebook",
					ID:       makeResponse.Value,
					Name:     name,
				}
				makes = append(makes, make)
			}
		}
	}

	p.logger.Debugf("%v", makes)

	return makes, nil
}

// GetModels gets models from facebook
func (p *Provider) GetModels(brand string) ([]structs.Model, error) {
	p.logger.Notice("GetModels for Facebook")

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	getModelsRequest := GetMakesRequest{
		CategoryIDS:            []int64{},
		CategoryRankingEnabled: true,
		ContextualData: []ContextualData{
			{Name: "seo_url", Value: "\"" + brand + "\""},
		},
		HideL2Cats: false,
		Surface:    "CATEGORY_FEED",
		TopicPageParams: TopicPageParams{
			LocationID: "category",
			URL:        brand,
		},
	}

	modelsRequestBody, err := json.Marshal(getModelsRequest)
	if err != nil {
		return []structs.Model{}, err
	}

	formData := url.Values{
		"fb_api_caller_class":      {"RelayModern"},
		"fb_api_req_friendly_name": {"CometMarketplaceSearchRootQuery"},
		"variables":                {string(modelsRequestBody)},
		"doc_id":                   {"5183985588339655"},
	}

	// p.logger.Debugf("%s", string(modelsRequestBody))

	req, err := http.NewRequest("POST", "https://en-gb.facebook.com/api/graphql/", strings.NewReader(formData.Encode()))
	if err != nil {
		return []structs.Model{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return []structs.Model{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []structs.Model{}, err
	}

	// p.logger.Debugf("%s", string(body))

	var modelsResponse GetMakesResponse
	if err := json.Unmarshal(body, &modelsResponse); err != nil {
		return []structs.Model{}, err
	}

	carModels := []structs.Model{}
	for _, filterType := range modelsResponse.Data.Viewer.MarketplaceStructuredFields {
		if filterType.FilterKey == "model" {
			for _, modelResponse := range filterType.Choices {
				carModel := structs.Model{
					Provider: "facebook",
					ID:       modelResponse.Value,
					Name:     modelResponse.DisplayName,
				}
				carModels = append(carModels, carModel)
			}
		}
	}

	// p.logger.Debugf("%v", carModels)

	return carModels, nil
}

func getLocation(postcode string) (structs.Location, error) {
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	locationReq, err := http.NewRequest("GET", "http://api.postcodes.io/postcodes/"+postcode, new(bytes.Buffer))
	if err != nil {
		return structs.Location{}, err
	}

	locationRes, err := client.Do(locationReq)
	if err != nil {
		return structs.Location{}, err
	}
	defer locationRes.Body.Close()

	locationBody, err := ioutil.ReadAll(locationRes.Body)
	if err != nil {
		return structs.Location{}, err
	}

	var locationResponse LocationResponse
	if err := json.Unmarshal(locationBody, &locationResponse); err != nil {
		return structs.Location{}, err
	}

	location := structs.Location{Latitude: locationResponse.Result.Latitude, Longitude: locationResponse.Result.Longitude}

	return location, nil
}

func getCityLocation(city string) (structs.Location, error) {
	db := database.DB

	if !strings.Contains(city, "United Kingdom") {
		city = city + ", United Kingdom"
	}

	var cityDB models.City
	err := db.Where("name = ?", city).First(&cityDB).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return structs.Location{}, err
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return structs.Location{Longitude: cityDB.Longitude, Latitude: cityDB.Latitude}, nil
	}

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	params := url.Values{}
	params.Add("key", config.Config("MAPQUEST_API_KEY"))
	params.Add("location", city)

	urlQuery := fmt.Sprintf("/geocoding/v1/address?%s", params.Encode())

	locationReq, err := http.NewRequest("GET", "http://open.mapquestapi.com"+urlQuery, new(bytes.Buffer))
	if err != nil {
		return structs.Location{}, err
	}

	locationRes, err := client.Do(locationReq)
	if err != nil {
		return structs.Location{}, err
	}
	defer locationRes.Body.Close()

	locationBody, err := ioutil.ReadAll(locationRes.Body)
	if err != nil {
		return structs.Location{}, err
	}

	var locationResponse LocationCityResponse
	if err := json.Unmarshal(locationBody, &locationResponse); err != nil {
		return structs.Location{}, err
	}

	longitude := locationResponse.Results[0].Locations[0].LatLng.Lng
	latitude := locationResponse.Results[0].Locations[0].LatLng.Lat

	cityDB = models.City{
		Name:      city,
		Longitude: longitude,
		Latitude:  latitude,
	}
	if err := db.Create(&cityDB).Error; err != nil {
		return structs.Location{}, err
	}

	location := structs.Location{Longitude: longitude, Latitude: latitude}

	return location, nil
}

func haversine(lonFrom float64, latFrom float64, lonTo float64, latTo float64) uint64 {
	var deltaLat = (latTo - latFrom) * (math.Pi / 180)
	var deltaLon = (lonTo - lonFrom) * (math.Pi / 180)

	var a = math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(latFrom*(math.Pi/180))*math.Cos(latTo*(math.Pi/180))*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := (6371 * c) + 1

	distanceString := uint64(distance) + 1
	return distanceString
}

// getAdvertsPaginationQuery gets specfic number adverts from index point from facebook
func getAdvertsPaginationQuery(getAdvertsPaginationQueryParams GetAdvertsPaginationQueryParams) ([]Edge, error) {
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	cursor := Cursor{
		Basic: Basic{
			ItemIndex: getAdvertsPaginationQueryParams.ItemIndex,
		},
	}

	cursorJSON, err := json.Marshal(cursor)
	if err != nil {
		return []Edge{}, err
	}

	getAdvertsRequest := GetAdvertsRequest{
		BuyLocation: BuyLocation{
			Latitude:  53.478698730469,
			Longitude: -2.2467041015625,
		},
		CategoryIDArray: []int64{807311116002614},
		ContextualData: []ContextualData{
			{Name: "make", Value: getAdvertsPaginationQueryParams.brand},
			{Name: "model", Value: getAdvertsPaginationQueryParams.model},
		},
		Count:                    getAdvertsPaginationQueryParams.Count,
		Cursor:                   string(cursorJSON),
		MarketplaceBrowseContext: "CATEGORY_FEED",
		PriceRange:               []int64{0, 214748364700},
		Radius:                   60000,
		Scale:                    2,
		StringVerticalFields: []ContextualData{
			{Name: "canonical_make_id", Value: getAdvertsPaginationQueryParams.brand},
			{Name: "canonical_model_id", Value: getAdvertsPaginationQueryParams.model},
		},
	}

	modelsRequestBody, err := json.Marshal(getAdvertsRequest)
	if err != nil {
		return []Edge{}, err
	}

	formData := url.Values{
		"fb_api_caller_class":      {"RelayModern"},
		"fb_api_req_friendly_name": {"CometMarketplaceCategoryContentPaginationQuery"},
		"variables":                {string(modelsRequestBody)},
		"doc_id":                   {"5111225608947754"},
	}

	// p.logger.Debugf("%s", string(modelsRequestBody))

	req, err := http.NewRequest("POST", "https://en-gb.facebook.com/api/graphql/", strings.NewReader(formData.Encode()))
	if err != nil {
		return []Edge{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return []Edge{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Edge{}, err
	}

	// p.logger.Debugf("%s", string(body))

	var advertsResponse GetAdvertsResponse
	if err := json.Unmarshal(body, &advertsResponse); err != nil {
		return []Edge{}, err
	}

	if len(advertsResponse.Data.Viewer.MarketplaceFeedStories.Edges) < 1 {
		return []Edge{}, nil
	}

	return advertsResponse.Data.Viewer.MarketplaceFeedStories.Edges, nil
}
