package model

type User struct {
	Id       int    `json:"id"`
	Username string `json:"nick"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
