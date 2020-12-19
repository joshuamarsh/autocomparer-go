package providers

// GetModelsParameters parser for query params
type GetModelsParameters struct {
	Provider []string `query:"provider[]"`
	Brand    string   `query:"make"`
}

// GetMakesParameters parser for query params
type GetMakesParameters struct {
	Provider []string `query:"provider[]"`
}

// GetAdvertsParameters parser for query params
type GetAdvertsParameters struct {
	Provider []string `query:"provider"`
	Brand    string   `query:"make"`
	Model    string   `query:"model"`
	Postcode string   `query:"postcode"`
	Radius   string   `query:"radius"`
	SortBy   string   `query:"sort"`
}
