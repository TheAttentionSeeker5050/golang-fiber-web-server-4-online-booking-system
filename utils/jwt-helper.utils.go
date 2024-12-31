package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func GenerateLocalAuthJWTToken(id string, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id, // We dont need to set an expiration time here because we will use the default expiration time from cookies. These will dissapear when the browser is closed or after the days in environment variable COOKIE_EXPIRES_IN_DAYS
		"email": email,
	})

	// get env secret in []byte
	var jwtSecret []byte = []byte(os.Getenv("JWT_SECRET"))
	if jwtSecret == nil {
		return "", errors.New("JWT_SECRET is not set")
	}

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return t, nil
}

func VerifyLocalAuthJWTToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

// The get token claims function
func GetLocalAuthJWTTokenClaims(tokenString string) (*JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	// if there are no email or id in the claims
	if _, ok := token.Claims.(jwt.MapClaims)["id"]; !ok {
		return nil, errors.New("id not in claims")
	}

	if _, ok := token.Claims.(jwt.MapClaims)["email"]; !ok {
		return nil, errors.New("email not in claims")
	}

	var mappedJWTClaims JWTClaims
	mappedJWTClaims.ID = token.Claims.(jwt.MapClaims)["id"].(string)
	mappedJWTClaims.Email = token.Claims.(jwt.MapClaims)["email"].(string)

	// if the id is not in the claims and bind it to the claims type struct
	// Return the reference to the struct mappedJWTClaims
	return &mappedJWTClaims, nil
}
