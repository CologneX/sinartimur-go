package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// ComparePasswords compares hashed password with plain password
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	// Logic to compare hashed password with plain password using bcrypt
	byteHash := []byte(hashedPwd)
	bytePlain := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err != nil {
		return false
	}
	return true
}

// HashPassword hashes plain password
func HashPassword(plainPwd string) string {
	// Logic to hash plain password using bcrypt
	bytePlain := []byte(plainPwd)
	hashedPwd, err := bcrypt.GenerateFromPassword(bytePlain, bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashedPwd)
}
