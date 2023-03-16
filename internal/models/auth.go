package models

type SignUpResponse struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required"`
	// TODO Avatar
}

type SignInResponse struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required"`
	// TODO Avatar
}

// TODO: these structures are the same. maybe rename and make one?
