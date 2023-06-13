package util

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPasword(plainPassword string) string {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), 14)
	return string(hashedPass)
}

func PasswordVerified(plainPassword, hashedPass string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPass),
		[]byte(plainPassword))
	if err != nil {
		return false
	}

	return true
}
