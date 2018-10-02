package userlogic

import (
	"strings"

	dbstore "github.com/kurt-midas/realworld-swagger/engines/datastore"
	sec "github.com/kurt-midas/realworld-swagger/engines/security"
)

type userLogicEngine struct {
	datastore dbstore.IDatastoreEngine
	security  sec.ISecurityEngine
}

func NewUserLogicEngine(db dbstore.IDatastoreEngine, sec sec.ISecurityEngine) IUserLogicEngine {
	return &userLogicEngine{datastore: db, security: sec}
}

func (engine *userLogicEngine) CreateUser(username, email, password string) (IUserEngine, error) {
	emailCipher, err := engine.security.Encrypt([]byte(email))
	if err != nil {
		return nil, err
	}
	salt, err := engine.security.NewSalt()
	if err != nil {
		return nil, err
	}

	pwhash := engine.security.Hash([]byte(password), salt)
	emailHash := engine.security.Hash([]byte(strings.ToLower(email)), nil)

	err = engine.datastore.NewUser(username, emailHash[:], emailCipher, pwhash[:], salt)
	if err != nil {
		return nil, err
	}
	return engine.GetUserByUsername(username)
}

func (engine *userLogicEngine) GetUserByEmail(email string) (IUserEngine, error) {
	emailHash := engine.security.Hash([]byte(strings.ToLower(email)), nil)
	dbuser, err := engine.datastore.GetUserByEmailHash(emailHash[:])
	if err != nil {
		return nil, err
	}
	var user = userEngine{
		datastore: engine.datastore,
		sec:       engine.security}
	user.info = dbuser
	m, err := engine.security.Decrypt(user.info.Email_crypt)
	if err != nil {
		return nil, err
	}
	user.email_plain = string(m)
	return &user, nil
}

func (engine *userLogicEngine) GetUserByUsername(username string) (IUserEngine, error) {
	dbuser, err := engine.datastore.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	var user = userEngine{
		datastore: engine.datastore,
		sec:       engine.security}
	user.info = dbuser
	m, err := engine.security.Decrypt(user.info.Email_crypt)
	if err != nil {
		return nil, err
	}
	user.email_plain = string(m)
	return &user, nil
}

func (engine *userLogicEngine) GetUserById(id int) (IUserEngine, error) {
	dbuser, err := engine.datastore.GetUserById(id)
	if err != nil {
		return nil, err
	}
	var user = userEngine{
		datastore: engine.datastore,
		sec:       engine.security}
	user.info = dbuser
	m, err := engine.security.Decrypt(user.info.Email_crypt)
	if err != nil {
		return nil, err
	}
	user.email_plain = string(m)
	return &user, nil
}
