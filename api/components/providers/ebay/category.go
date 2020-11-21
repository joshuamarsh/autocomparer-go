package ebay

import "encoding/xml"

type GetCategoryInfoResponse struct {
	XMLName       xml.Name `xml:"GetCategoryInfoResponse"`
	Text          string   `xml:",chardata"`
	Xmlns         string   `xml:"xmlns,attr"`
	Timestamp     string   `xml:"Timestamp"`
	Ack           string   `xml:"Ack"`
	Build         string   `xml:"Build"`
	Version       string   `xml:"Version"`
	CategoryArray struct {
		Text     string `xml:",chardata"`
		Category []struct {
			Text             string `xml:",chardata"`
			CategoryID       string `xml:"CategoryID"`
			CategoryLevel    string `xml:"CategoryLevel"`
			CategoryName     string `xml:"CategoryName"`
			CategoryParentID string `xml:"CategoryParentID"`
			CategoryNamePath string `xml:"CategoryNamePath"`
			CategoryIDPath   string `xml:"CategoryIDPath"`
			LeafCategory     string `xml:"LeafCategory"`
		} `xml:"Category"`
	} `xml:"CategoryArray"`
	CategoryCount   string `xml:"CategoryCount"`
	UpdateTime      string `xml:"UpdateTime"`
	CategoryVersion string `xml:"CategoryVersion"`
}
