package database

type User struct {
	ID    int
	Email string
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUSer(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	ID := len(dbStructure.Users) + 1
	user := User{
		ID:    ID,
		Email: email,
	}

	dbStructure.Users[ID] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
