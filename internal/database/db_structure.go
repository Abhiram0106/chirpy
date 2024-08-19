package database

type DBStructure struct {
	Chirps map[int]Chirp        `json:"chirps"`
	Users  map[int]internalUser `json:"users"`
}
