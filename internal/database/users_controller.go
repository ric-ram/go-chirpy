package database

import (
	"errors"
)

type User struct {
	ID         int
	Email      string
	Password   string
	IsChirpRed bool
}

var ErrUserAlreadyExists = errors.New("user already exists")

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUSer(email, password string) (User, error) {
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrUserAlreadyExists
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	ID := len(dbStructure.Users) + 1
	user := User{
		ID:         ID,
		Email:      email,
		Password:   password,
		IsChirpRed: false,
	}

	dbStructure.Users[ID] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetUserByEmail returns the user with the corresponded email
func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}

// GetUserByID returns the user with the corresponded id
func (db *DB) GetUserByID(userID int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[userID]
	if !ok {
		return User{}, ErrNotExist
	}

	return user, nil
}

// UpdateUser returns the updated user
func (db *DB) UpdateUser(id int, email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, err := db.GetUserByID(id)
	if err != nil {
		return User{}, ErrNotExist
	}

	user.Email = email
	user.Password = password
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// UpgradeUser returns the upgraded user to chirpy red
func (db *DB) UpgradeUser(user User) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user.IsChirpRed = true
	dbStructure.Users[user.ID] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
