package models

// Make struct
type Make struct {
	Value         string `gorm:"primary_key;not null"`
	Provider      string `gorm:"primary_key;not null"`
	Name          string `gorm:"not null"`
	ProviderValue string `gorm:"not null"`
}
