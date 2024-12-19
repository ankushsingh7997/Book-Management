package routes

import (
	"github.com/ankush/bookstore/pkg/controllers"
	"github.com/ankush/bookstore/pkg/middlewares"
	"github.com/gorilla/mux"
)

var RegisterBookRoutes = func(router *mux.Router) {
	router.HandleFunc("/books", middlewares.LogRequestResponse(controllers.CreateBook)).Methods("POST")
	router.HandleFunc("/book/", middlewares.LogRequestResponse(controllers.GetBook)).Methods("GET")
	router.HandleFunc("/book/{bookId}", middlewares.LogRequestResponse(controllers.GetBookById)).Methods("GET")
	router.HandleFunc("/book/{bookId}", middlewares.LogRequestResponse(controllers.UpdateBook)).Methods("PUT")
	router.HandleFunc("book/{bookId}", middlewares.LogRequestResponse(controllers.DeleteBook)).Methods("DELETE")

}
