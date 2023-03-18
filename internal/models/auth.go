package models

type SignUpResponse struct {
	Email string `json:"email"`
	Image Image  `json:"image"`
}

type SignInResponse struct {
	Email string `json:"email"`
	Image Image  `json:"image"`
}
