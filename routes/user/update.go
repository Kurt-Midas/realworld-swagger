package user

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/kurt-midas/realworld-swagger/models"
)

func (u *UserApiEngine) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	authheader := r.Header.Get("Authorization")
	tokens := strings.Split(authheader, " ")
	if len(tokens) != 2 {
		log.Print("Unrecognized token format\n")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Printf("DecodeJWT(%s)\n", tokens[1])
	id, err := u.SessionEngine.DecodeUserToken(tokens[1])
	if err != nil {
		log.Printf("Update User token decode failed: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var rbody models.UpdateUserRequest
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&rbody)

	user, err := u.UserLogic.GetUserById(id)
	if err != nil {
		log.Printf("Failed to retrieve user with id %d: %s\n", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rbody.User.Email != "" {
		if _, err := mail.ParseAddress(rbody.User.Email); err != nil {
			log.Printf("Invalid new mail: %s\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if len(rbody.User.Password) > 0 && len(rbody.User.Password) < 8 {
		log.Print("Invalid new password\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = user.UpdateUser(*rbody.User)
	if err != nil {
		log.Printf("Failed to update user id %d: %s\n", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// refresh token, return user
	info := user.GetInfo()
	info.Token, err = u.SessionEngine.CreateUserToken(info.Id)
	if err != nil {
		log.Printf("Update user successful but could not generate new token: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := models.UserResponse{}
	response.User = &info
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&response)
}
