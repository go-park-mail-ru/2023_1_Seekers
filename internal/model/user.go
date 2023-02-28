package model

type User struct {
	Id       uint64 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Profile struct {
	UId       uint64 `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
}

type FormSignUp struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	RepeatPw  string `json:"repeat_pw"` // ?
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
}

// ? omit

type FormLogin struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	BirthDate string `json:"birth_date"`
}
