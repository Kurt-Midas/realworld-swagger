package models

type User struct {
	Id int `json:"-"`

	Email string `json:"email"`

	Token string `json:"token"`

	Username string `json:"username"`

	Bio *string `json:"bio"`

	Image *string `json:"image"`
}
