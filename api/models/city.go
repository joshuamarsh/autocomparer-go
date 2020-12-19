package models

// City struct
type City struct {
	Name      string  `gorm:"primary_key;not null"`
	Longitude float64 `gorm:"not null"`
	Latitude  float64 `gorm:"not null"`
}
