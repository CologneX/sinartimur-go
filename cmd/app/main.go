package main

import (
	"database/sql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	v1 "sinartimur-go/api/v1"
	"sinartimur-go/config"
	"sinartimur-go/internal/auth"
	"sinartimur-go/utils"
)

type Services struct {
	UserService *auth.AuthService
}

func main() {
	// start a connection to Postgres
	db := config.StartPostgres()
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Failed to close database: %v", err)
		}
	}()
	// start a connection to Redis
	redisClient := config.NewRedisClient()
	// register services
	services := registerServices(db, redisClient)

	// define routes
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	v1.RegisterUserRoutes(router, services.UserService)

	// Add logging middleware from gorilla/mux
	loggedRouter := handlers.CustomLoggingHandler(os.Stdout, router, utils.Logger)

	// serve the router on port 8080
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(loggedRouter)))
}

func registerServices(
	db *sql.DB,
	redis *config.RedisClient,
) *Services {
	userRepo := auth.NewUserRepository(db)
	userService := auth.NewAuthService(userRepo, redis)
	// Initialize other services here

	return &Services{
		UserService: userService,
	}
}
