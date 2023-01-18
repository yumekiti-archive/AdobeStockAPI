package config

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"AdobeStockAPI/domain"
)

func NewDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database")
	}

	db.AutoMigrate(&domain.Body{})

	return db
}
