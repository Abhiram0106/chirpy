package database

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	JWT          string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
}
