package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("./ideb.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	// Migrate the schema
	if err := DB.AutoMigrate(&Request{}, &CorporateRequest{}, &GetIdeb{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
