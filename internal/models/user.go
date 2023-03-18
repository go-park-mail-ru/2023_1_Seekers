package models

type User struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Avatar    string `json:"avatar,omitempty"`
}

type FormSignUp struct {
	Login     string `json:"login" validate:"required"`
	Password  string `json:"password" validate:"required"`
	RepeatPw  string `json:"repeat_pw" validate:"required"` // ?
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type FormLogin struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	Remember bool   `json:"remember" validate:"required"`
}
