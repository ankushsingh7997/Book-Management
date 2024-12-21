package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ankush/bookstore/pkg/models"
	"github.com/ankush/bookstore/pkg/utils"
	"github.com/gorilla/mux"
)

var NewBook models.Book

func CreateBook(res http.ResponseWriter, req *http.Request) {
	CreateBook := &models.Book{}
	utils.ParseBody(req, CreateBook)

	b := CreateBook.CreateBook()
	resp, _ := json.Marshal(b)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)

}
func GetBook(res http.ResponseWriter, req *http.Request) {
	newBook := models.GetAllBooks()
	data, _ := json.Marshal(newBook)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(data)

}

func GetBookById(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	BookId := vars["bookId"]
	ID, err := strconv.ParseInt(BookId, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing")
	}
	result, _, err := models.GetBookById(ID)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(`{"status:false,"message":"Not Found"}`))
		return

	}
	data, _ := json.Marshal(result)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(data)

}
func UpdateBook(res http.ResponseWriter, req *http.Request) {
	UpdateBook := &models.Book{}
	utils.ParseBody(req, UpdateBook)
	vars := mux.Vars(req)
	BookId := vars["bookId"]
	ID, err := strconv.ParseInt(BookId, 0, 0)
	if err != nil {
		fmt.Println("Wrror while parsing")
	}
	BookDetails, db, _ := models.GetBookById(ID)
	if UpdateBook.Name != "" {
		BookDetails.Name = UpdateBook.Name
	}
	if UpdateBook.Author != "" {
		BookDetails.Author = UpdateBook.Author
	}
	if UpdateBook.Publication != "" {
		BookDetails.Publication = UpdateBook.Publication
	}
	db.Save(&BookDetails)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(BookDetails)
	res.Write(data)

}
func DeleteBook(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	BookId := vars["bookId"]
	ID, err := strconv.ParseInt(BookId, 0, 0)
	if err != nil {
		log.Panic(err)
	}
	data := models.DeleteBook(ID)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(data)
	res.Write(resp)

}
