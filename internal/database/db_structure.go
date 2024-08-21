package database

import "time"

type timeAndID struct {
	TokenExpiresAt time.Time
	UserID         int
}

type DBStructure struct {
	NextChirpID   int                  `json:"next_chirp_id"`
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]internalUser `json:"users"`
	RefreshTokens map[string]timeAndID `json:"refresh_tokens"`
}

func (dbStruct *DBStructure) FindUserByEmail(email string) (internalUser, bool) {

	for _, user := range dbStruct.Users {
		if user.Email == email {
			return user, true
		}
	}

	return internalUser{}, false
}
