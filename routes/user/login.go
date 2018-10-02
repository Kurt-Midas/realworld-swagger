package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kurt-midas/realworld-swagger/models"
)

// Login - Existing user login
func (u *UserApiEngine) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var body models.LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Error in Login json decode: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if body.User.Email == "" || len(body.User.Password) < 8 {
		log.Print("Error in Login: invalid required fields\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := u.UserLogic.GetUserByEmail(body.User.Email)
	if err != nil {
		log.Printf("Error in Login get user: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !user.IsValidPassword(body.User.Password) {
		log.Print("Error in Login: invalid password\n")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	info := user.GetInfo()
	info.Token, err = u.SessionEngine.CreateUserToken(info.Id)
	if err != nil {
		log.Printf("Login failed at token generation: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := models.UserResponse{}
	response.User = &info
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&response)
}
