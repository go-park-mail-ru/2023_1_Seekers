package model

type FormReg struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	RepeatPw string `json:"repeat_pw"`
}
