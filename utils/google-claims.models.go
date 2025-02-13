package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoogleClaims struct {
	ID            primitive.ObjectID `json:"id"`
	Email         string             `json:"email"`
	EmailVerified bool               `json:"email_verified"`
	Sub           string             `json:"sub"` // Google user ID
	Name          string             `json:"name"`
	Picture       string             `json:"picture"`
}

func VerifyGoogleToken(token string) (*GoogleClaims, error) {
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var claims GoogleClaims
	if err := json.Unmarshal(body, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
