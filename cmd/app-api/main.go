package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	v1 "sinartimur-go/api/v1"
	"sinartimur-go/config"
	"sinartimur-go/internal/user"
	"sinartimur-go/utils"
)

type Services struct {
	UserService *user.UserService
}

func main() {
	// start a connection to Postgres
	db := config.StartPostgres()
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Failed to close database: %v", err)
		}
	}()

	// register services
	services := registerServices()

	// define routes
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	v1.RegisterUserRoutes(router, services.UserService)

	// Add logging middleware from gorilla/mux
	loggedRouter := handlers.CustomLoggingHandler(os.Stdout, router, utils.Logger)

	// serve the router on port 8080
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(loggedRouter)))
}

func registerServices() *Services {
	userRepo := user.NewUserRepository()
	userService := user.NewUserService(userRepo)
	// Initialize other services here

	return &Services{
		UserService: userService,
	}
}
