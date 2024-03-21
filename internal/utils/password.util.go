package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SALT_SIZE  = 32
	HASH_SIZE  = 64
	ITERATIONS = 10000
	SEPARATOR  = "."
)

func HashPassword(password string) (string, error) {
	saltBytes := make([]byte, SALT_SIZE)

	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", err
	}

	hashedPasswordBytes := pbkdf2.Key(
		[]byte(password),
		saltBytes,
		ITERATIONS,
		HASH_SIZE,
		sha256.New,
	)

	salt := base64.StdEncoding.Strict().EncodeToString(saltBytes)
	hashedPassword := base64.StdEncoding.Strict().EncodeToString(hashedPasswordBytes)

	hash := hashedPassword + SEPARATOR + salt
	return hash, nil
}

func VerifyPassword(password, hashedPassword string) (bool, error) {
	hashes := strings.Split(hashedPassword, SEPARATOR)
	if len(hashes) < 2 {
		return false, &CustomError{
			Message: "Invalid password",
		}
	}

	hashedPasswordBytes, err := base64.StdEncoding.DecodeString(hashes[0])
	if err != nil {
		return false, err
	}

	salt, err := base64.StdEncoding.DecodeString(hashes[1])
	if err != nil {
		return false, err
	}

	newHashedPassword := pbkdf2.Key([]byte(password), salt, ITERATIONS, HASH_SIZE, sha256.New)

	if subtle.ConstantTimeCompare(hashedPasswordBytes, newHashedPassword) != 1 {
		return false, nil
	}

	return true, nil
}
