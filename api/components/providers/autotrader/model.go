package autotrader

// GetModelsRequest json object request
type GetModelsRequest struct {
	OperationName string    `json:"operationName"`
	Variables     Variables `json:"variables"`
	Query         string    `json:"query"`
}

// AdvertQuery reponse object parser from autotrader
type AdvertQuery struct {
	Make []string `json:"make"`
}

// Variables reponse object parser from autotrader
type Variables struct {
	AdvertQuery AdvertQuery `json:"advertQuery"`
	Facets      []string    `json:"facets"`
}

// GetModelsResponse reponse object parser from autotrader
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

// Values reponse object parser from autotrader
type Values struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Facets reponse object parser from autotrader
type Facets struct {
	Name   string   `json:"name"`
	Values []Values `json:"values"`
}
