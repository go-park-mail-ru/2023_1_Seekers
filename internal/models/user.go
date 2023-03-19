package models

type User struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Avatar    string `json:"avatar,omitempty"`
}

type FormSignUp struct {
	Login     string `json:"login" validate:"required"`
	Password  string `json:"password" validate:"required"`
	RepeatPw  string `json:"repeatPw" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type FormLogin struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	Remember bool   `json:"remember" validate:"required"`
}
