package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ankush/bookstore/logger"
	"github.com/ankush/bookstore/pkg/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const port string = ":3001"

func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Unhandled error: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Close server gracefully
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a new router using Gorilla Mux
	r := mux.NewRouter()
	r.Use(ErrorHandlerMiddleware)

	// Register the routes for users and books
	apiRouter := r.PathPrefix("/api/v1").Subrouter() // Subrouter for versioning
	routes.RegisterBookRoutes(apiRouter.PathPrefix("/book").Subrouter())
	routes.RegisterUserRoutes(apiRouter.PathPrefix("/user").Subrouter())

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
		logg := logger.NewLogger("UserService", "production")
		logg.Info(fmt.Sprintf("Listening on port: %s", port))

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
