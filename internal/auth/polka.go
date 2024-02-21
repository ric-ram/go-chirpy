package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GetPolkaApiKey get's the polka api key from the header
func GetPolkaApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

// ValidatePolkaApiKey validates polka api key
func ValidatePolkaApiKey(headerApiKey string, polkaApiSecret string) error {
	if headerApiKey != polkaApiSecret {
		return fmt.Errorf("invalid polka api key")
	}

	return nil
}
