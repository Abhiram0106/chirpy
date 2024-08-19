package database

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"sync"
)

type internalUser struct {
	ID       int
	Email    string
	Password []byte
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

func NewDB(path string) (*DB, error) {

	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {

	database, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	newChirpId := len(database.Chirps) + 1

	newChirp := Chirp{
		ID:   newChirpId,
		Body: body,
	}

	database.Chirps[newChirpId] = newChirp

	writeDBError := db.writeDB(database)

	if writeDBError != nil {
		return Chirp{}, writeDBError
	}

	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {

	database, err := db.loadDB()

	if err != nil {
		return nil, err
	}

	chirps := []Chirp{}

	for _, chirp := range database.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) CreateUser(email string, password string) (User, error) {

	database, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	newUserId := len(database.Users) + 1

	hashedPassword, hashingErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if hashingErr != nil {
		return User{}, hashingErr
	}

	newInternalUser := internalUser{
		ID:       newUserId,
		Email:    email,
		Password: hashedPassword,
	}

	database.Users[newUserId] = newInternalUser

	writeDBError := db.writeDB(database)

	if writeDBError != nil {
		return User{}, writeDBError
	}

	newUser := User{
		ID:    newUserId,
		Email: email,
	}

	return newUser, nil
}

func (db *DB) GetUsers() ([]User, error) {

	database, err := db.loadDB()

	if err != nil {
		return nil, err
	}

	users := []User{}

	for _, user := range database.Users {
		users = append(users, User{
			ID:    user.ID,
			Email: user.Email,
		})
	}

	return users, nil
}

func (db *DB) ensureDB() error {

	defer db.mux.Unlock()
	db.mux.Lock()

	_, statErr := os.Stat(db.path)

	if statErr == nil {
		log.Println("DATABASE ALREADY EXISTS")
		return nil
	}

	if !(errors.Is(statErr, os.ErrNotExist)) {
		log.Printf("ensureDB 1 %s\n", statErr.Error())
		return statErr
	}

	log.Println("DATABASE DOESN'T EXIST")

	emptyDB := DBStructure{
		Chirps: make(map[int]Chirp),
		Users:  make(map[int]internalUser),
	}

	dbJson, marshalErr := json.Marshal(emptyDB)

	if marshalErr != nil {
		return marshalErr
	}

	writeErr := os.WriteFile(db.path, dbJson, 0666)
	if writeErr != nil {
		return writeErr
	}

	log.Println("DATABASE CREATED")

	return nil
}

func (db *DB) loadDB() (DBStructure, error) {

	defer db.mux.Unlock()
	db.mux.Lock()

	data, readErr := os.ReadFile(db.path)

	if readErr != nil {
		return DBStructure{}, readErr
	}

	database := DBStructure{}
	marshalErr := json.Unmarshal(data, &database)
	if marshalErr != nil {
		return DBStructure{}, marshalErr
	}

	log.Printf("DB LOADED %v", database)
	return database, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {

	defer db.mux.Unlock()
	db.mux.Lock()

	body, marshalErr := json.Marshal(dbStructure)

	if marshalErr != nil {
		return marshalErr
	}

	err := os.WriteFile(db.path, body, 0666)

	if err != nil {
		return err
	}
	log.Println("DB WRITTEN")

	return nil
}
