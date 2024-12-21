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

func (b *Book) CreateBook() *Book {
	db.NewRecord(b)
	db.Create(&b)
	return b

}

func GetAllBooks() []Book {
	var Books []Book
	db.Find(&Books)
	return Books
}

func GetBookById(id int64) (*Book, *gorm.DB, error) {
	var getBook Book
	result := db.Where("ID=?", id).First(&getBook)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	return &getBook, result, nil
}

func DeleteBook(id int64) Book {
	var book Book
	db.Where("ID=?", id).Delete(&book)
	return book
}
