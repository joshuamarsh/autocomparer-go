package database

import (
	"carcompare/api/config"
	"carcompare/api/models"
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql database driver
)

// DB gorm connector
var DB *gorm.DB

// Connect connects to db
func Connect() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)

	DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_HOST"), port, config.Config("DB_NAME")))

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
	DB.AutoMigrate(&models.User{})
	fmt.Println("Database Migrated")
}
