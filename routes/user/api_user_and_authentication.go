package user

import (
	"github.com/kurt-midas/realworld-swagger/engines/session"
	"github.com/kurt-midas/realworld-swagger/engines/userlogic"
)

/*
type IUserApiEngine interface {
	GetCurrentUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateCurrentUser(w http.ResponseWriter, r *http.Request)
} */

type UserApiEngine struct {
	UserLogic     userlogic.IUserLogicEngine
	SessionEngine session.ISessionEngine
}

func NewUserApiEngine(usrlogic userlogic.IUserLogicEngine, sess session.ISessionEngine) *UserApiEngine {
	return &UserApiEngine{
		UserLogic:     usrlogic,
		SessionEngine: sess}
}
