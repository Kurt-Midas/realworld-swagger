package user

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"

	"github.com/kurt-midas/realworld-swagger/models"
)

// CreateUser - Register a new user
func (u *UserApiEngine) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var body models.NewUserRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Error 400 decoding json request: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if body.User.Username == "" {
		log.Printf("Error 400: empty username\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m, err := mail.ParseAddress(body.User.Email)
	if err != nil {
		log.Printf("Error 400 at parsing mail: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(body.User.Password) < 8 {
		log.Printf("Error 400, password too short\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := u.UserLogic.CreateUser(body.User.Username, m.Address, body.User.Password)
	if err != nil {
		log.Printf("Could not create user: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	info := user.GetInfo()
	info.Token, err = u.SessionEngine.CreateUserToken(info.Id)
	if err != nil {
		log.Printf("User created but token generation failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := models.UserResponse{}
	response.User = &info
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&response)
}
