package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type internalUser struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    []byte `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
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

func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {

	database, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	newChirpId := database.NextChirpID
	database.NextChirpID++

	newChirp := Chirp{
		ID:       newChirpId,
		Body:     body,
		AuthorID: authorID,
	}

	database.Chirps[newChirpId] = newChirp

	if writeDBError := db.writeDB(database); writeDBError != nil {
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

func (db *DB) GetChirpByID(chirpID int) (Chirp, error) {

	database, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	chirp, exists := database.Chirps[chirpID]
	if !exists {
		return Chirp{}, errors.New("Chirp not found")
	}

	return chirp, nil
}

func (db *DB) DeleteChirpByID(chirpID int, authorID int) error {

	database, err := db.loadDB()

	if err != nil {
		return err
	}

	chirp, exists := database.Chirps[chirpID]
	if !exists {
		return errors.New("Chirp not found")
	}

	if chirp.AuthorID != authorID {
		return errors.New("Unauthorized")
	}

	delete(database.Chirps, chirpID)

	if writeDBErr := db.writeDB(database); writeDBErr != nil {
		return writeDBErr
	}

	return nil
}

func (db *DB) CreateUser(email string, password string) (User, error) {

	database, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	if _, exists := database.FindUserByEmail(email); exists {
		return User{}, errors.New("Email in use")
	}

	newUserId := len(database.Users) + 1

	hashedPassword, hashingErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if hashingErr != nil {
		return User{}, hashingErr
	}

	newInternalUser := internalUser{
		ID:          newUserId,
		Email:       email,
		Password:    hashedPassword,
		IsChirpyRed: false,
	}

	database.Users[newUserId] = newInternalUser

	if writeDBError := db.writeDB(database); writeDBError != nil {
		return User{}, writeDBError
	}

	newUser := User{
		ID:          newUserId,
		Email:       email,
		IsChirpyRed: false,
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
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		})
	}

	return users, nil
}

func (db *DB) GetUserByEmailAndPassword(email string, password string) (User, error) {

	database, err := db.loadDB()

	if err != nil {
		return User{}, err
	}

	interUser, exists := database.FindUserByEmail(email)

	if !exists {
		return User{}, errors.New("Invalid email or password")
	}

	if passwordErr := bcrypt.CompareHashAndPassword(interUser.Password, []byte(password)); passwordErr != nil {
		return User{}, errors.New("Invalid email or password")
	}

	user := User{
		ID:          interUser.ID,
		Email:       interUser.Email,
		IsChirpyRed: interUser.IsChirpyRed,
	}

	return user, nil
}

func (db *DB) UpdateUser(userID int, email string, password string) (User, error) {

	database, loadDBErr := db.loadDB()

	if loadDBErr != nil {
		return User{}, loadDBErr
	}

	user := database.Users[userID]

	if _, exists := database.FindUserByEmail(email); exists && user.Email != email {
		return User{}, errors.New("Email in use")
	}

	if len(email) != 0 {
		user.Email = email
	}

	if len(password) != 0 {
		hashedPassword, hashingErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if hashingErr != nil {
			return User{}, hashingErr
		}
		user.Password = hashedPassword
	}

	database.Users[userID] = user

	writeErr := db.writeDB(database)

	if writeErr != nil {
		return User{}, writeErr
	}

	updatedUser := User{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	return updatedUser, nil
}

func (db *DB) AddRefreshToken(token string, expires_at time.Time, userID int) error {

	database, loadDBErr := db.loadDB()

	if loadDBErr != nil {
		return loadDBErr
	}

	database.RefreshTokens[token] = timeAndID{TokenExpiresAt: expires_at, UserID: userID}

	dbWriteErr := db.writeDB(database)

	if dbWriteErr != nil {
		return dbWriteErr
	}

	return nil
}

func (db *DB) IsRefreshTokenValid(token string) (userID int, err error) {

	database, loadDBErr := db.loadDB()

	if loadDBErr != nil {
		return 0, loadDBErr
	}

	tokenMetaData, exists := database.RefreshTokens[token]
	if !exists {
		return 0, errors.New("Token doesn't exist")
	}

	if tokenMetaData.TokenExpiresAt.Before(time.Now().UTC()) {
		return 0, errors.New("Token has expired")
	}

	internalUser, exists := database.Users[tokenMetaData.UserID]

	if !exists {
		return 0, errors.New("User doesn't exist")
	}

	return internalUser.ID, nil
}

func (db *DB) RevokeRefreshToken(token string) error {

	database, loadDBErr := db.loadDB()

	if loadDBErr != nil {
		return loadDBErr
	}

	if _, exists := database.RefreshTokens[token]; !exists {
		return errors.New("Token doesn't exist")
	}
	delete(database.RefreshTokens, token)

	dbWriteErr := db.writeDB(database)

	if dbWriteErr != nil {
		return dbWriteErr
	}

	return nil
}

func (db *DB) UpgradeUserToChirpyRed(userID int) error {

	database, loadDBErr := db.loadDB()

	if loadDBErr != nil {
		return loadDBErr
	}

	user, exists := database.Users[userID]

	if !exists {
		return errors.New("User not found")
	}

	user.IsChirpyRed = true

	database.Users[userID] = user

	writeDBErr := db.writeDB(database)

	if writeDBErr != nil {
		return writeDBErr
	}

	return nil
}

func (db *DB) ensureDB() error {

	defer db.mux.Unlock()
	db.mux.Lock()

	_, statErr := os.Stat(db.path)

	if statErr == nil {
		return nil
	}

	if !(errors.Is(statErr, os.ErrNotExist)) {
		log.Printf("ensureDB 1 %s\n", statErr.Error())
		return statErr
	}

	emptyDB := DBStructure{
		Chirps:        make(map[int]Chirp),
		Users:         make(map[int]internalUser),
		RefreshTokens: make(map[string]timeAndID),
		NextChirpID:   1,
	}

	dbJson, marshalErr := json.Marshal(emptyDB)

	if marshalErr != nil {
		return marshalErr
	}

	writeErr := os.WriteFile(db.path, dbJson, 0666)
	if writeErr != nil {
		return writeErr
	}

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

	return nil
}
