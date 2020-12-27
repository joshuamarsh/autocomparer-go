package ebay

import (
	"bytes"
	"carcompare/config"
	"carcompare/structs"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/apsdehal/go-logger"
	"github.com/leekchan/accounting"
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

// GetAdvert gets adverts from ebay
func (p *Provider) GetAdvert(postcode string, radius string, brand string, model string, sortBy string, page *uint) ([]structs.Adverts, error) {
	p.logger.Notice("GetAdvert for Ebay")

	ebaySort := ""
	switch sortBy {
	case "best_match":
		ebaySort = "BestMatch"
	case "date_desc":
		ebaySort = "StartTimeNewest"
	case "date_asc":
		ebaySort = ""
	case "dist_asc":
		ebaySort = "DistanceNearest"
	case "year_desc":
		ebaySort = "year-desc"
	case "price_asc":
		ebaySort = "PricePlusShippingLowest"
	case "price_desc":
		ebaySort = "PricePlusShippingHighest"
	case "miles_asc":
		ebaySort = ""
	default:
		ebaySort = ""
	}

	params := url.Values{}
	params.Add("OPERATION-NAME", "findItemsAdvanced")
	params.Add("paginationInput.entriesPerPage", "30")
	params.Add("keywords", brand+" "+model)
	params.Add("buyerPostalCode", postcode)
	params.Add("outputSelector", "PictureURLLarge")
	params.Add("sortOrder", ebaySort)
	params.Add("categoryId", "9801")
	params.Add("itemFilter(0).name", "MaxDistance")
	params.Add("itemFilter(0).value", radius)

	if page == nil || *page == 0 || *page == 1 {
		params.Add("paginationInput.pageNumber", "1")
	} else {
		ebayPage := strconv.Itoa(int(*page))
		if ebayPage == "" {
			params.Add("paginationInput.pageNumber", "1")
		}
		params.Add("paginationInput.pageNumber", ebayPage)
	}

	urlQuery := fmt.Sprintf("/services/search/FindingService/v1?%s", params.Encode())

	body, err := p.execute("GET", urlQuery)
	if err != nil {
		return []structs.Adverts{}, err
	}

	var advertsResponse FindItemsByKeywordsResponse
	if err := xml.Unmarshal(body, &advertsResponse); err != nil {
		return []structs.Adverts{}, err
	}

	adverts := structs.Adverts{}
	for _, advertResponse := range advertsResponse.SearchResult.Item {
		advertImage := advertResponse.PictureURLLarge
		if advertImage == "" {
			advertImage = advertResponse.GalleryURL
		}

		priceFormatted := ""
		if priceFloat, err := strconv.ParseFloat(advertResponse.SellingStatus.CurrentPrice.Text, 64); err == nil {
			ac := accounting.Accounting{Symbol: "£", Precision: 2, Thousand: ",", Decimal: "."}
			priceFormatted = ac.FormatMoney(priceFloat)
		}
		priceReplace := strings.NewReplacer("£", "", ",", "", ".", "")
		priceFloat := priceReplace.Replace(priceFormatted)
		pricePence, _ := strconv.ParseInt(priceFloat, 10, 64)

		location := strings.ReplaceAll(advertResponse.Location, ",", ", ")
		distanceFloat, _ := strconv.ParseFloat(advertResponse.Distance.Text, 64)

		advert := structs.Advert{
			Provider:    "ebay",
			ID:          advertResponse.ItemID,
			Link:        advertResponse.ViewItemURL,
			Location:    location,
			Distance:    uint64(distanceFloat),
			Title:       advertResponse.Title,
			Price:       pricePence,
			Description: advertResponse.Subtitle,
			Image:       advertImage,
		}
		adverts.AddAdvert(advert)
	}

	providerAdverts := make([]structs.Adverts, 0)
	providerAdverts = append(providerAdverts, adverts)
	return providerAdverts, nil
}

func (p *Provider) execute(method string, parameters string) ([]byte, error) {

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	buffer := new(bytes.Buffer)

	req, err := http.NewRequest(method, "https://svcs.ebay.com"+parameters, buffer)
	if err != nil {
		p.logger.Fatalf("%s", err)
	}

	req.Header.Add("X-EBAY-SOA-SECURITY-APPNAME", config.Config("EBAY_APP_ID"))
	req.Header.Add("X-EBAY-SOA-GLOBAL-ID", "EBAY-GB")
	req.Header.Add("X-EBAY-SOA-RESPONSE-DATA-FORMAT", "XML")
	req.Header.Add("X-EBAY-SOA-SERVICE-VERSION", "1.0.0")

	res, err := client.Do(req)
	if err != nil {
		p.logger.Fatalf("%s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		p.logger.Fatalf("%s", err)
	}

	// p.logger.DebugF(" %s response from ebay: %s", method, html.UnescapeString(string(body)))

	return body, nil
}

// GetMakes gets adverts from ebay
func (p *Provider) GetMakes() ([]structs.Make, error) {
	p.logger.Notice("GetMakes for Ebay")

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	buffer := new(bytes.Buffer)

	req, err := http.NewRequest("GET", "http://open.api.ebay.com/Shopping?callname=GetCategoryInfo&CategoryID=9801&IncludeSelector=ChildCategories", buffer)
	if err != nil {
		return []structs.Make{}, err
	}

	req.Header.Add("X-EBAY-API-APP-NAME", config.Config("EBAY_APP_ID"))
	req.Header.Add("X-EBAY-API-SITE-ID", "3")
	req.Header.Add("X-EBAY-API-CALL-NAME", "GetCategoryInfo")
	req.Header.Add("X-EBAY-API-VERSION", "1157")

	res, err := client.Do(req)
	if err != nil {
		return []structs.Make{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []structs.Make{}, err
	}

	var categoriesResponse GetCategoryInfoResponse
	if err := xml.Unmarshal(body, &categoriesResponse); err != nil {
		return []structs.Make{}, err
	}

	categories := []structs.Make{}
	for _, categoryResponse := range categoriesResponse.CategoryArray.Category {
		if categoryResponse.CategoryLevel == "3" {
			name := categoryResponse.CategoryName
			switch name {
			case "Citroën":
				name = "Citroen"
			case "Seat":
				name = "SEAT"
			case "Rolls Royce":
				name = "Rolls-Royce"
			case "Rover/MG":
				advert := structs.Make{
					Provider: "ebay",
					ID:       "18262",
					Name:     "MG",
				}
				categories = append(categories, advert)
				name = "Rover"
			case "Kit Cars":
				name = "Kit Cars/Custom Vehicle"
			case "Other Cars":
				name = ""
			}
			if name != "" && categoryResponse.CategoryID != "" {
				advert := structs.Make{
					Provider: "ebay",
					ID:       categoryResponse.CategoryID,
					Name:     name,
				}
				categories = append(categories, advert)
			}
		}
	}

	return categories, nil
}

// GetModels gets models from ebay
func (p *Provider) GetModels(brand string) ([]structs.Model, error) {
	p.logger.Notice("GetModels for Ebay")
	return []structs.Model{}, nil
}
