package userlogic

import (
	"bytes"
	"strings"

	dbstore "github.com/kurt-midas/realworld-swagger/engines/datastore"
	"github.com/kurt-midas/realworld-swagger/engines/security"
	"github.com/kurt-midas/realworld-swagger/models"
)

type userEngine struct {
	sec         security.ISecurityEngine
	datastore   dbstore.IDatastoreEngine
	info        dbstore.DB_user
	email_plain string
}

func (u *userEngine) GetInfo() models.User {
	info := models.User{}
	info.Id = u.info.Id
	info.Username = u.info.Username
	info.Bio = u.info.Bio
	info.Image = u.info.Image
	info.Email = u.email_plain
	// info.Token = u.info.Token
	return info
}

func (u *userEngine) IsValidPassword(password string) bool {
	pw := u.sec.Hash([]byte(password), u.info.Password_salt)
	return bytes.Equal(pw[:], u.info.Password_hash)
}

func (u *userEngine) UpdateUser(newInfo models.UpdateUser) error {
	info := u.GetInfo()
	tempinfo := u.info
	var isNewEmail bool
	updater := u.datastore.UpdateUser(info.Id)
	if info.Username != newInfo.Username {
		updater = updater.Username(newInfo.Username)
		tempinfo.Username = newInfo.Username
	}
	if info.Bio != newInfo.Bio {
		updater = updater.Bio(newInfo.Bio)
		tempinfo.Bio = newInfo.Bio
	}
	if info.Image != newInfo.Image {
		updater = updater.Image(newInfo.Image)
		tempinfo.Image = newInfo.Image
	}
	if info.Email != newInfo.Email {
		ehash := u.sec.Hash([]byte(strings.ToLower(newInfo.Email)), nil)
		ecrypt, err := u.sec.Encrypt([]byte(newInfo.Email))
		if err != nil {
			return err
		}
		updater = updater.Email_hash(ehash[:]).Email_crypt(ecrypt)
		tempinfo.Email_crypt = ecrypt
		tempinfo.Email_hash = ehash[:]
		isNewEmail = true
	}
	if !u.IsValidPassword(newInfo.Password) {
		salt, err := u.sec.NewSalt()
		if err != nil {
			return err
		}
		phash := u.sec.Hash([]byte(newInfo.Password), salt)
		updater = updater.Password_hash(phash[:]).Password_salt(salt)
		tempinfo.Password_hash = phash[:]
		tempinfo.Password_salt = salt
	}
	err := updater.Execute()
	if err != nil {
		return err
	}
	// no errors. Update the info.
	u.info = tempinfo
	if isNewEmail {
		u.email_plain = newInfo.Email
	}
	return nil
}
