package routes

import (
	"github.com/ankush/bookstore/pkg/controllers"
	"github.com/ankush/bookstore/pkg/middlewares"
	"github.com/gorilla/mux"
)

var RegisterBookRoutes = func(router *mux.Router) {
	router.HandleFunc("/create", middlewares.LogRequestResponse(controllers.CreateBook)).Methods("POST")
	router.HandleFunc("/getAll", middlewares.LogRequestResponse(controllers.GetBook)).Methods("GET")
	router.HandleFunc("/get/{bookId}", middlewares.LogRequestResponse(controllers.GetBookById)).Methods("GET")
	router.HandleFunc("/update/{bookId}", middlewares.LogRequestResponse(controllers.UpdateBook)).Methods("PUT")
	router.HandleFunc("/delete/{bookId}", middlewares.LogRequestResponse(controllers.DeleteBook)).Methods("DELETE")

}
