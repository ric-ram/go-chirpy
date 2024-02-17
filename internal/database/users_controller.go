package database

type User struct {
	ID       int
	Email    string
	Password string
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUSer(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	_, exists, err := db.GetUserByEmail(email)
	if err != nil || exists {
		return User{}, err
	}

	ID := len(dbStructure.Users) + 1
	user := User{
		ID:       ID,
		Email:    email,
		Password: password,
	}

	dbStructure.Users[ID] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetUserByEmail returns the user with the corresponded email
func (db *DB) GetUserByEmail(email string) (User, bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, false, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, true, nil
		}
	}

	return User{}, false, nil
}
