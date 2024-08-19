package database

type DBStructure struct {
	Chirps map[int]Chirp        `json:"chirps"`
	Users  map[int]internalUser `json:"users"`
}

func (dbStruct *DBStructure) FindUserByEmail(email string) (internalUser, bool) {

	for _, user := range dbStruct.Users {
		if user.Email == email {
			return user, true
		}
	}

	return internalUser{}, false
}
