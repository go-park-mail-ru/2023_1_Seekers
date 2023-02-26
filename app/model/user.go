package model

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Profile struct {
	UId       int    `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
}

type FormSignUp struct {
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	RepeatPw  string `json:"repeat_pw,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	BirthDate string `json:"birth_date,omitempty"`
}

type FormLogin struct {
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	BirthDate string `json:"birth_date,omitempty"`
}
