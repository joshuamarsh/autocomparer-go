package providers

type GetModelsParameters struct {
	Provider []string `query:"provider[]"`
	Brand    string   `query:"make"`
}

type GetMakesParameters struct {
	Provider []string `query:"provider[]"`
}

type GetAdvertsParameters struct {
	Provider []string `query:"provider"`
	Brand    string   `query:"make"`
	Model    string   `query:"model"`
	Postcode string   `query:"postcode"`
	Radius   string   `query:"radius"`
	SortBy   string   `query:"sort"`
}
