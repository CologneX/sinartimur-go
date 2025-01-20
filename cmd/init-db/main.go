package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sinartimur-go/internal/user"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <username> <password>\n", os.Args[0])
	}

	username := os.Args[1]
	password := os.Args[2]

	conn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	defer db.Close()

	// Create user service
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)

	// Create user
	req := user.CreateUserRequest{
		Username:        username,
		Password:        password,
		ConfirmPassword: password,
	}

	httpCode, err := userService.CreateUser(req)
	if err != nil {
		log.Fatalf("Failed to create user: %v (HTTP code: %d)\n", err, httpCode)
	}

	fmt.Println("User successfully created")
}
