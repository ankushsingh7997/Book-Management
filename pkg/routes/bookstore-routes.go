package routes

import (
	"github.com/ankush/bookstore/pkg/controllers"
	"github.com/ankush/bookstore/pkg/middlewares"
	"github.com/gorilla/mux"
)

var RegisterBookRoutes = func(router *mux.Router) {
	protect := router.PathPrefix("").Subrouter()
	protect.Use(middlewares.Protect)
	protect.HandleFunc("/create", middlewares.LogRequestResponse(controllers.CreateBook)).Methods("POST")
	protect.HandleFunc("/getAll", middlewares.LogRequestResponse(controllers.GetBook)).Methods("GET")
	protect.HandleFunc("/get/{bookId}", middlewares.LogRequestResponse(controllers.GetBookById)).Methods("GET")
	protect.HandleFunc("/update/{bookId}", middlewares.LogRequestResponse(controllers.UpdateBook)).Methods("PUT")
	protect.HandleFunc("/delete/{bookId}", middlewares.LogRequestResponse(controllers.DeleteBook)).Methods("DELETE")

}
