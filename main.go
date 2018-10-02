package main

import (
	"log"
	"net/http"

	"encoding/hex"

	"github.com/kurt-midas/realworld-swagger/engines/datastore"
	"github.com/kurt-midas/realworld-swagger/engines/security"
	"github.com/kurt-midas/realworld-swagger/engines/userlogic"
	"github.com/kurt-midas/realworld-swagger/engines/session"
	"github.com/kurt-midas/realworld-swagger/routes"
	"github.com/kurt-midas/realworld-swagger/routes/user"
)

var aeskey [32]byte
var salt []byte
var secret []byte
var sessionExpiration = 1800 // 30 minutes

func init() {
	key, _ := hex.DecodeString("abcdffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	copy(aeskey[:], key)
	salt = []byte("default salt")
	secret, _ = hex.DecodeString("session signing secret")
}

func main() {
	log.Printf("Server started")

	db, err := datastore.NewMySQL("root", "NblCMfOTyz2DMv49aI1z", "127.0.0.1:3306", "realworld")
	if err != nil {
		log.Printf("db err, go figure: %s\n", err)
	}
	sec := security.NewSecurityEngine(aeskey, salt)
	usrlogic := userlogic.NewUserLogicEngine(db, sec)
	sess := session.NewJwtSessionEngine(secret, sessionExpiration)

	usrapi := user.NewUserApiEngine(usrlogic, sess)
	r := routes.MyRouter{U: usrapi}

	router := r.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
