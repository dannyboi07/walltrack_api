package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPw), err
}

func VerifyPassword(hashPw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPw), []byte(pw)) == nil
}
