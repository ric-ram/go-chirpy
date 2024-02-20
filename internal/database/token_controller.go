package database

import (
	"errors"
	"time"
)

type RevokedToken struct {
	ID         string
	RevokeTime time.Time
}

var ErrTokenAlreadyExists = errors.New("the refesh token is already revoked")

// AddRevokeToken adds the refresh token as revoked to the database
func (db *DB) AddRevokeToken(token string) error {
	if _, err := db.GetRevokedTokenById(token); !errors.Is(err, ErrNotExist) {
		return ErrTokenAlreadyExists
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	ID := token
	revokedToken := RevokedToken{
		ID:         token,
		RevokeTime: time.Now().UTC(),
	}

	dbStructure.RevokedTokens[ID] = revokedToken

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

// GetRevokedTokenById returns the revoked token with the specified ID
func (db *DB) GetRevokedTokenById(tokenID string) (RevokedToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RevokedToken{}, err
	}

	revokedToken, ok := dbStructure.RevokedTokens[tokenID]
	if !ok {
		return RevokedToken{}, ErrNotExist
	}

	return revokedToken, nil
}
