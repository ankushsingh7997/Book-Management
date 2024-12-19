package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ankush/bookstore/pkg/routes"
	"github.com/gorilla/mux"
)

const port string = ":3001"

func main() {
	// Close server gracefully
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// Create a new router using Gorilla Mux
	r := mux.NewRouter()

	// Register the routes for users and books
	apiRouter := r.PathPrefix("/api/v1").Subrouter() // Subrouter for versioning
	routes.RegisterBookRoutes(apiRouter.PathPrefix("/books").Subrouter())

	// Initialize the HTTP server with the router
	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a separate goroutine
	go func() {
		log.Println("Listening on port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error while connecting to server: %v\n", err)
		}
	}()

	// Wait for shutdown signal
	<-stopChan
	log.Println("Shutting down server")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
