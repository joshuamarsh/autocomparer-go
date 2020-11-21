package autotrader

type GetModelsRequest struct {
	OperationName string    `json:"operationName"`
	Variables     Variables `json:"variables"`
	Query         string    `json:"query"`
}
type AdvertQuery struct {
	Make []string `json:"make"`
}
type Variables struct {
	AdvertQuery AdvertQuery `json:"advertQuery"`
	Facets      []string    `json:"facets"`
}

type GetModelsResponse struct {
	Data struct {
		Search struct {
			Adverts struct {
				Facets []struct {
					Name   string   `json:"name"`
					Values []Values `json:"values"`
				} `json:"facets"`
			} `json:"adverts"`
		} `json:"search"`
	} `json:"data"`
}

type Values struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Facets struct {
	Name   string   `json:"name"`
	Values []Values `json:"values"`
}
