package main

import (
	"log"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitDatabase(dbPath string) (*gorm.DB, error) {
	// Get the absolute path of the database file
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		log.Printf("Error getting absolute path for database: %v", err)
	} else {
		log.Printf("Using database file at: %s", absPath)
	}

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