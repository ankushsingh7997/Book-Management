package models

import (
	"github.com/ankush/bookstore/pkg/config"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// Book represents the book model
type Book struct {
	gorm.Model         // Correctly embed gorm.Model for fields like ID, CreatedAt, etc.
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

// Initialize the `db` variable
func init() {
	config.Connect() // Ensure the DB connection is established
	db = config.GetDB()
	db.AutoMigrate(&Book{}) // Automatically migrate the schema
}
