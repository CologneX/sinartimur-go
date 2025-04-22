package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
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

	// hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Gagal Hash Password: %v", err)
	}

	// Create user in the database with admin role
	_, err = db.Exec("Insert Into Appuser (Username, Password_Hash, Is_Admin) Values ($1, $2, $3)", username, string(passwordHash), true)
	if err != nil {
		log.Fatalf("Gagal membuat user: %v", err)
	}

	fmt.Println("User berhasil dibuat")
}
