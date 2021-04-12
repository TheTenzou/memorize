package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

func hashPassword(password string) (string, error) {
	salt := make([]byte, 32)

	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	hashedPassword := fmt.Sprintf(
		"%s.%s",
		hex.EncodeToString(hash),
		hex.EncodeToString(salt),
	)

	return hashedPassword, nil
}

func comparePassword(storedPassword, suppliedPassword string) (bool, error) {
	passwordAndSalt := strings.Split(storedPassword, ".")

	salt, err := hex.DecodeString(passwordAndSalt[1])

	if err != nil {
		return false, fmt.Errorf("unable to verify user password")
	}

	hash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)

	if err != nil {
		return false, fmt.Errorf("unable to verify user password")
	}

	return hex.EncodeToString(hash) == passwordAndSalt[0], nil
}
