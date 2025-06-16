package services

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	nickname string `'json:"nickname"`
	Password string `json:"password"`
}