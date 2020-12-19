package ebay

import (
	"encoding/xml"
)

// FindItemsByKeywordsResponse reponse parser from ebay
type FindItemsByKeywordsResponse struct {
	XMLName      xml.Name `xml:"findItemsAdvancedResponse"`
	Text         string   `xml:",chardata"`
	Xmlns        string   `xml:"xmlns,attr"`
	Ack          string   `xml:"ack"`
	Version      string   `xml:"version"`
	Timestamp    string   `xml:"timestamp"`
	SearchResult struct {
		Text  string `xml:",chardata"`
		Count string `xml:"count,attr"`
		Item  []struct {
			Text            string `xml:",chardata"`
			ItemID          string `xml:"itemID"`
			Title           string `xml:"title"`
			GlobalID        string `xml:"globalID"`
			Subtitle        string `xml:"subtitle"`
			PrimaryCategory struct {
				Text         string `xml:",chardata"`
				CategoryID   string `xml:"categoryID"`
				CategoryName string `xml:"categoryName"`
			} `xml:"primaryCategory"`
			GalleryURL      string   `xml:"galleryURL"`
			PictureURLLarge string   `xml:"pictureURLLarge"`
			ViewItemURL     string   `xml:"viewItemURL"`
			PaymentMethod   []string `xml:"paymentMethod"`
			AutoPay         string   `xml:"autoPay"`
			PostalCode      string   `xml:"postalCode"`
			Location        string   `xml:"location"`
			Country         string   `xml:"country"`
			ShippingInfo    struct {
				Text                string `xml:",chardata"`
				ShippingServiceCost struct {
					Text       string `xml:",chardata"`
					CurrencyID string `xml:"currencyID,attr"`
				} `xml:"shippingServiceCost"`
				ShippingType    string `xml:"shippingType"`
				ShipToLocations string `xml:"shipToLocations"`
			} `xml:"shippingInfo"`
			SellingStatus struct {
				Text         string `xml:",chardata"`
				CurrentPrice struct {
					Text       string `xml:",chardata"`
					CurrencyID string `xml:"currencyID,attr"`
				} `xml:"currentPrice"`
				ConvertedCurrentPrice struct {
					Text       string `xml:",chardata"`
					CurrencyID string `xml:"currencyID,attr"`
				} `xml:"convertedCurrentPrice"`
				SellingState string `xml:"sellingState"`
				TimeLeft     string `xml:"timeLeft"`
				BIDCount     string `xml:"bIDCount"`
			} `xml:"sellingStatus"`
			ListingInfo struct {
				Text              string `xml:",chardata"`
				BestOfferEnabled  string `xml:"bestOfferEnabled"`
				BuyItNowAvailable string `xml:"buyItNowAvailable"`
				StartTime         string `xml:"startTime"`
				EndTime           string `xml:"endTime"`
				ListingType       string `xml:"listingType"`
				Gift              string `xml:"gift"`
				WatchCount        string `xml:"watchCount"`
			} `xml:"listingInfo"`
			Distance struct {
				Text string `xml:",chardata"`
				Unit string `xml:"unit,attr"`
			} `xml:"distance"`
			Condition struct {
				Text                 string `xml:",chardata"`
				ConditionID          string `xml:"conditionID"`
				ConditionDisplayName string `xml:"conditionDisplayName"`
			} `xml:"condition"`
			IsMultiVariationListing string `xml:"isMultiVariationListing"`
			DiscountPriceInfo       struct {
				Text                string `xml:",chardata"`
				OriginalRetailPrice struct {
					Text       string `xml:",chardata"`
					CurrencyID string `xml:"currencyID,attr"`
				} `xml:"originalRetailPrice"`
				PricingTreatment string `xml:"pricingTreatment"`
				SoldOnEbay       string `xml:"soldOnEbay"`
				SoldOffEbay      string `xml:"soldOffEbay"`
			} `xml:"discountPriceInfo"`
			TopRatedListing string `xml:"topRatedListing"`
			ProductID       struct {
				Text string `xml:",chardata"`
				Type string `xml:"type,attr"`
			} `xml:"productID"`
			SecondaryCategory struct {
				Text         string `xml:",chardata"`
				CategoryID   string `xml:"categoryID"`
				CategoryName string `xml:"categoryName"`
			} `xml:"secondaryCategory"`
		} `xml:"item"`
	} `xml:"searchResult"`
	PaginationOutput struct {
		Text           string `xml:",chardata"`
		PageNumber     string `xml:"pageNumber"`
		EntriesPerPage string `xml:"entriesPerPage"`
		TotalPages     string `xml:"totalPages"`
		TotalEntries   string `xml:"totalEntries"`
	} `xml:"paginationOutput"`
	ItemSearchURL string `xml:"itemSearchURL"`
}
