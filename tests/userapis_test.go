package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kurt-midas/realworld-swagger/engines/datastore"
	"github.com/kurt-midas/realworld-swagger/engines/security"
	"github.com/kurt-midas/realworld-swagger/engines/session"
	"github.com/kurt-midas/realworld-swagger/engines/userlogic"
	"github.com/kurt-midas/realworld-swagger/models"
	"github.com/kurt-midas/realworld-swagger/routes/user"
)

func Test_UserApiIntegration(t *testing.T) {
	var mockkey [32]byte
	_, err := rand.Read(mockkey[:])
	if err != nil {
		t.Fatalf("Failed to create mock key: %s\n", err)
	}
	salt := []byte("mock testing salt")
	secret := []byte("mock testing secret")

	db, err := datastore.NewMySQL("root", "NblCMfOTyz2DMv49aI1z", "127.0.0.1:3306", "realworld")
	if err != nil {
		t.Fatalf("db err, go figure: %s\n", err)
	}
	sec := security.NewSecurityEngine(mockkey, salt)
	usrlogic := userlogic.NewUserLogicEngine(db, sec)
	sess := session.NewJwtSessionEngine(secret, 60)

	usrapi := user.NewUserApiEngine(usrlogic, sess)

	utf := time.Now().Unix()
	username := fmt.Sprintf("Test_User_%d", utf)
	email := fmt.Sprintf("Test_User_%d@localhost", utf)
	password := "Password"
	authtoken := ""

	// Create User
	{
		var createbuff bytes.Buffer
		createbody := models.NewUserRequest{
			User: &models.NewUser{
				Username: username,
				Email:    email,
				Password: password,
			},
		}
		err = json.NewEncoder(&createbuff).Encode(&createbody)
		if err != nil {
			t.Fatalf("Failed to encode create user request body: %s\n", err)
		}

		ts := httptest.NewServer(http.HandlerFunc(usrapi.CreateUser))

		req, _ := http.NewRequest(http.MethodPost, ts.URL, &createbuff)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Create User request failed %s\n", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected response 201 created, got %d\n", resp.StatusCode)
		}
		defer resp.Body.Close()
		var respbody models.UserResponse
		err = json.NewDecoder(resp.Body).Decode(&respbody)
		if respbody.User.Email != email {
			t.Errorf("Login response expected email %s, got %s\n", email, respbody.User.Email)
		}
		if respbody.User.Username != username {
			t.Errorf("Login response expected username %s, got %s\n", username, respbody.User.Username)
		}
		if respbody.User.Token == "" {
			t.Error("Create User response expected a token but got none\n")
		}
	}

	// Login
	{
		var loginbuff bytes.Buffer
		loginbody := models.LoginUserRequest{
			User: &models.LoginUser{
				Email:    email,
				Password: password,
			},
		}
		err = json.NewEncoder(&loginbuff).Encode(&loginbody)
		if err != nil {
			t.Fatalf("Failed to encode login request body: %s\n", err)
		}

		ts := httptest.NewServer(http.HandlerFunc(usrapi.Login))

		req, _ := http.NewRequest(http.MethodPost, ts.URL, &loginbuff)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Login request failed %s\n", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected response 200 OK, got %d\n", resp.StatusCode)
		}
		defer resp.Body.Close()
		var respbody models.UserResponse
		err = json.NewDecoder(resp.Body).Decode(&respbody)
		if err != nil {
			t.Fatalf("Failed to decode login response: %s\n", err)
		}
		if respbody.User.Email != email {
			t.Errorf("Login response expected email %s, got %s\n", email, respbody.User.Email)
		}
		authtoken = respbody.User.Token
	}

	// Get Current User
	{
		ts := httptest.NewServer(http.HandlerFunc(usrapi.GetCurrentUser))

		req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
		req.Header.Add("Authorization", "Token "+authtoken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Current User request failed %s\n", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected response 200 OK, got %d\n", resp.StatusCode)
		}
		defer resp.Body.Close()
		var respbody models.UserResponse
		err = json.NewDecoder(resp.Body).Decode(&respbody)
		if err != nil {
			t.Fatalf("Failed to decode current user response: %s\n", err)
		}
		if respbody.User.Email != email {
			t.Errorf("Current User response expected email %s, got %s\n", email, respbody.User.Email)
		}
	}

	var newEmail = "Updated_" + email
	var newPassword = "Updated_" + password
	var newUsername = "Updated_" + username

	// Update
	{
		var updatebuff bytes.Buffer
		updatebody := models.UpdateUserRequest{
			User: &models.UpdateUser{
				Email:    newEmail,
				Password: newPassword,
				Username: newUsername,
			},
		}
		err = json.NewEncoder(&updatebuff).Encode(&updatebody)
		if err != nil {
			t.Fatalf("Failed to encode Update request body: %s\n", err)
		}

		ts := httptest.NewServer(http.HandlerFunc(usrapi.UpdateCurrentUser))

		req, _ := http.NewRequest(http.MethodPost, ts.URL, &updatebuff)
		req.Header.Add("Authorization", "Token "+authtoken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Update request failed %s\n", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected response 200 OK, got %d\n", resp.StatusCode)
		}
		defer resp.Body.Close()
		var respbody models.UserResponse
		err = json.NewDecoder(resp.Body).Decode(&respbody)
		if err != nil {
			t.Fatalf("Failed to decode update response: %s\n", err)
		}
		if respbody.User.Email != newEmail {
			t.Errorf("Update response expected newEmail %s, got %s\n", newEmail, respbody.User.Email)
		}
		if respbody.User.Username != newUsername {
			t.Errorf("Update response expected newUsername %s, got %s\n", newUsername, respbody.User.Username)
		}
	}

	// Fail login
	{
		var loginbuff bytes.Buffer
		loginbody := models.LoginUserRequest{
			User: &models.LoginUser{
				Email:    newEmail,
				Password: password, //old password should be invalid
			},
		}
		err = json.NewEncoder(&loginbuff).Encode(&loginbody)
		if err != nil {
			t.Fatalf("Failed to encode login fail case request body: %s\n", err)
		}

		ts := httptest.NewServer(http.HandlerFunc(usrapi.Login))

		req, _ := http.NewRequest(http.MethodPost, ts.URL, &loginbuff)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Login fail case request failed %s\n", err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("Expected response 401 Unauthorized, got %d\n", resp.StatusCode)
		}
	}
}
