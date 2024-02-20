package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ErrNoAuthHeaderIncluded -
var ErrNoAuthHeaderIncluded = errors.New("not auth header included in request")

// HashPassword creates a hash password to safely store it in the database
func HashPassword(password string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(encryptedPassword), nil
}

// ValidatePassword validates the login password with the hash password of the user
func ValidatePassword(hashPassword, loginPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(loginPassword))
}

// GetBearerToken get's the token from the header
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

// CreateJwtToken creates the jwt token
func CreateJwtToken(userID int, tokenSecret string, jwtTokenType string) (string, error) {
	expiresIn, err := GetExpirationTime(jwtTokenType)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    fmt.Sprintf("chirpy-%s", jwtTokenType),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   fmt.Sprintf("%d", userID),
	})

	return token.SignedString([]byte(tokenSecret))
}

// ValidateJwtToken validates token
func ValidateJwtToken(headerToken string, tokenSecret string) (*jwt.Token, error) {
	claimStruct := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(headerToken, &claimStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// ValidateAccessJwtToken validates if the token is a valid jwt token
func ValidateAccessJwtToken(headerToken string, tokenSecret string) (*jwt.Token, error) {
	validToken, err := ValidateJwtToken(headerToken, tokenSecret)
	if err != nil {
		return nil, err
	}
	issuer, err := validToken.Claims.GetIssuer()
	if err != nil {
		return nil, err
	}

	isRefreshToken, err := IsRefreshToken(issuer)
	if err != nil {
		return nil, err
	}
	if isRefreshToken {
		return nil, fmt.Errorf("token is not a valid access token")
	}

	return validToken, nil
}

// GetUserID returns the userID present in the jwt Token
func GetUserID(token *jwt.Token) (string, error) {
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	return userIDString, nil
}

// ValidateRefreshJwtToken validates if the token is a valid jwt token
func ValidateRefreshJwtToken(headerToken string, tokenSecret string) (*jwt.Token, error) {
	validToken, err := ValidateJwtToken(headerToken, tokenSecret)
	if err != nil {
		return nil, err
	}

	issuer, err := validToken.Claims.GetIssuer()
	if err != nil {
		return nil, err
	}

	isRefreshToken, err := IsRefreshToken(issuer)
	if err != nil {
		return nil, err
	}
	if !isRefreshToken {
		return nil, fmt.Errorf("token is not a valid refresh token")
	}

	return validToken, nil
}

// GetExpirationTime returns the jwt token expiration time based on the token type
func GetExpirationTime(jwtTokenType string) (time.Duration, error) {
	accessExspiration := 60 * 60           // 1 hour in seconds
	refreshExpiration := 60 * 60 * 60 * 24 // 60 days in seconds
	if jwtTokenType == "access" {
		return time.Duration(accessExspiration) * time.Second, nil
	} else if jwtTokenType == "refresh" {
		return time.Duration(refreshExpiration) * time.Second, nil
	}

	return time.Duration(0), fmt.Errorf("invalid token type: %s", jwtTokenType)
}

// IsRefreshToken validates if the token is a refresh token
func IsRefreshToken(issuer string) (bool, error) {
	if issuer == "chirpy-refresh" {
		return true, nil
	} else if issuer == "chirpy-access" {
		return false, nil
	}

	return false, fmt.Errorf("not a valid issuer")
}
