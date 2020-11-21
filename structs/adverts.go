package structs

type AdvertProviders struct {
	Providers []string `json:"providers"`
	Adverts   []Advert `json:"adverts"`
}

type Adverts struct {
	Adverts []Advert `json:"adverts"`
}

type Advert struct {
	Provider    string `json:"provider"`
	ID          string `json:"id"`
	Link        string `json:"link"`
	Location    string `json:"location"`
	Distance    string `json:"distance"`
	Title       string `json:"title"`
	Price       int64  `json:"price"`
	Mileage     string `json:"mileage"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func (adverts *Adverts) AddAdvert(advert Advert) []Advert {
	adverts.Adverts = append(adverts.Adverts, advert)
	return adverts.Adverts
}

type Make struct {
	Provider string `json:"provider"`
	ID       string `json:"id"`
	Name     string `json:"name"`
}

type MakeProvider struct {
	ID        string            `json:"id"`
	Providers map[string]string `json:"providers"`
	Name      string            `json:"name"`
}

type Model struct {
	Provider string `json:"provider"`
	ID       string `json:"id"`
	Name     string `json:"name"`
}

type ModelProvider struct {
	ID        string            `json:"id"`
	Providers map[string]string `json:"providers"`
	Name      string            `json:"name"`
}
