package models

type UpdateUser struct {
	Email string `json:"email,omitempty"`

	Username string `json:"username,omitempty"`

	Bio *string `json:"bio,omitempty"`

	Image *string `json:"image,omitempty"`

	Password string `json:"password,omitempty"`
}
