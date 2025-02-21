package generatetoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT generates a JWT token using the provided secret key.
func GenerateJWT(secretKey string) (string, error) {
	// Define token claims
	claims := jwt.MapClaims{
		"authorized": true,
		"user":       "example_user",
		"exp":        time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
