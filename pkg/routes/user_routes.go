package routes

import (
	"github.com/ankush/bookstore/pkg/controllers"
	"github.com/ankush/bookstore/pkg/middlewares"
	"github.com/gorilla/mux"
)

// add protect middleware to the routes
var RegisterUserRoutes = func(router *mux.Router) {
	public := router.PathPrefix("").Subrouter()
	public.HandleFunc("/create", middlewares.LogRequestResponse(controllers.CreateUser)).Methods("POST")
	public.HandleFunc("/login", middlewares.LogRequestResponse(controllers.Login)).Methods("POST")
	protected := router.PathPrefix("").Subrouter()
	protected.Use(middlewares.Protect)
	protected.HandleFunc("/getAll", middlewares.LogRequestResponse(controllers.GetAllUsers)).Methods("GET")
	protected.HandleFunc("/get/{userId}", middlewares.LogRequestResponse(controllers.GetUserByID)).Methods("GET")
	protected.HandleFunc("/update/{userId}", middlewares.LogRequestResponse(controllers.UpdateUser)).Methods("PUT")
	protected.HandleFunc("/delete/{userId}", middlewares.LogRequestResponse(controllers.DeleteUser)).Methods("DELETE")

}
