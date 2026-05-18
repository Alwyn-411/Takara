package handlers

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CheckPassword(Hashed string, Plaintext string) error {
	return bcrypt.CompareHashAndPassword([]byte(Hashed), []byte(Plaintext))
}