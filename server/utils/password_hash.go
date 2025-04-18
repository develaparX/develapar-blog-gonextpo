package utils

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func ConparePasswordHash(passwordHash string, passwordFromDb string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordFromDb))
}
