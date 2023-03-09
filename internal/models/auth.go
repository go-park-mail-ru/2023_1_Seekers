package models

type SignUpResponse struct {
	Email string `json:"email" validate:"required"`
	// TODO Avatar
}

type SignInResponse struct {
	Email string `json:"email" validate:"required"`
	// TODO Avatar
}
