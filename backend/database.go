package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Request{}, &CorporateRequest{}, &GetIdeb{}); err != nil {
		return nil, err
	}

	return db, nil
}