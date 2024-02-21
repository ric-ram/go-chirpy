package database

import "errors"

var ErrNotExist = errors.New("resource does not exist")

type Chirp struct {
	ID       int
	Body     string
	AuthorID int
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	lenChirps := len(dbStructure.Chirps)
	if lenChirps == 0 {
		lenChirps = 1
	}
	ID := dbStructure.Chirps[lenChirps-1].ID + 1
	chirp := Chirp{
		ID:       ID,
		Body:     body,
		AuthorID: authorID,
	}

	dbStructure.Chirps[ID] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// GetChirpsById returns the chirp with the correspondent ID
func (db *DB) GetChirpsById(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

// DeleteChirp deletes the chirp
func (db *DB) DeleteChirp(chirp Chirp) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// delete(map, key)
	delete(dbStructure.Chirps, chirp.ID)

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
