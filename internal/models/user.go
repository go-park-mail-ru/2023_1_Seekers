package models

type User struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Profile struct {
	UID       uint64 `json:"user_id" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type FormSignUp struct {
	Email     string `json:"email" validate:"email,required"`
	Password  string `json:"password" validate:"required"`
	RepeatPw  string `json:"repeat_pw" validate:"required"` // ?
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type FormLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
