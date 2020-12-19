package models

// Model struct
type Model struct {
	Value         string `gorm:"primary_key;not null"`
	Provider      string `gorm:"primary_key;not null"`
	Make          string `gorm:"primary_key;not null"`
	Name          string `gorm:"not null"`
	ProviderValue string `gorm:"not null"`
}
