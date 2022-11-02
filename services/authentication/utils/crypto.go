package utils

import "golang.org/x/crypto/bcrypt"

//bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//bcrypt.GenerateFromPassword([]byte(password), 14)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
