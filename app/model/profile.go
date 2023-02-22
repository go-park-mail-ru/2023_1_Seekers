package model

type Profile struct {
	UId       int    `json:"user_id"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	NickName  string `json:"nick_name"`
	Sex       string `json:"sex"`
	City      string `json:"city"`
	Avatar    string `json:"avatar"`
}
