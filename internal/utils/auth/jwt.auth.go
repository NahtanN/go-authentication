package auth_utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Tokens struct {
	AccessToken            string    `json:"accessToken"  example:"eyJhbGciOiJIUzI1N..."`
	RefreshToken           string    `json:"refreshToken" example:"eyJhbGciOiJIUzI1N..."`
	RefreshTokenExpiration time.Time `json:"-"`
}

func CreateJwtTokens(id uint32) (*Tokens, error) {
	secret := os.Getenv("JWT_SECRET")
	byteSecret := []byte(secret)

	accessTokenClaims := jwt.MapClaims{
		"exp": time.Now().Add(10 * time.Minute).Unix(), // Expires in 10 minutes
		"id":  id,
	}
	accessTokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	refreshTokenExpiration := time.Now().AddDate(0, 0, 15)

	refreshTokenClaims := jwt.MapClaims{
		"exp": refreshTokenExpiration.Unix(),
	}
	refreshTokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	accessToken, err := accessTokenString.SignedString(byteSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := refreshTokenString.SignedString(byteSecret)
	if err != nil {
		return nil, err
	}

	tokens := Tokens{
		AccessToken:            accessToken,
		RefreshToken:           refreshToken,
		RefreshTokenExpiration: refreshTokenExpiration,
	}

	return &tokens, nil
}
