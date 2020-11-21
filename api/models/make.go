package models

import "github.com/jinzhu/gorm"

// Make struct
type Make struct {
	gorm.Model
	Value    string `gorm:"not null" json:"value"`
	Name     string `gorm:"not null" json:"password"`
	Provider string `gorm:"not null" json:"provider"`
}
