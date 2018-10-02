package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kurt-midas/realworld-swagger/models"
)

// GetCurrentUser - Get current user
func (u *UserApiEngine) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	authheader := r.Header.Get("Authorization")
	tokens := strings.Split(authheader, " ")
	if len(tokens) != 2 {
		log.Print("Unrecognized auth header format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("DecodeJWT(%s)\n", tokens[1])
	id, err := u.SessionEngine.DecodeUserToken(tokens[1])
	if err != nil {
		log.Printf("Get Current User token decode failed: %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := u.UserLogic.GetUserById(id)
	if err != nil {
		log.Printf("Failed to retrieve user with id %d: %s\n", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// refresh token, return user
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
