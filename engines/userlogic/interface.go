package userlogic

import (
	"github.com/kurt-midas/realworld-swagger/models"
)

type IUserLogicEngine interface {
	CreateUser(username, email, password string) (IUserEngine, error)
	GetUserByEmail(string) (IUserEngine, error)
	GetUserByUsername(string) (IUserEngine, error)
	GetUserById(int) (IUserEngine, error)
}

type IUserEngine interface {
	GetInfo() models.User

	IsValidPassword(string) bool

	UpdateUser(newInfo models.UpdateUser) error
}
