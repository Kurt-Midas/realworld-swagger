package datastore

import (
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type mysqlEngine struct {
	db *sqlx.DB
}

func NewMySQL(user string, pass string, address string, db_name string) (IDatastoreEngine, error) {
	var format string = "%s:%s@tcp(%s)/%s?charset=utf8"
	var datasource = fmt.Sprintf(format, user, pass, address, db_name)
	db, err := sqlx.Open("mysql", datasource)
	return &mysqlEngine{db}, err
}

func (m *mysqlEngine) NewUser(username string, email_hash, email_crypt, password_hash, password_salt []byte) error {
	var query = "INSERT INTO users (username, email_hash, email_crypt, password_hash, password_salt) VALUES (?, ?, ?, ?, ?)"
	args := []interface{}{username, email_hash, email_crypt, password_hash, password_salt}
	_, err := m.db.Exec(query, args...)
	// 1062 duplicate entry?
	return err
}

func (m *mysqlEngine) GetUserByEmailHash(emailhash []byte) (DB_user, error) {
	var query = "SELECT id, username, email_hash, email_crypt, password_hash, password_salt, bio, image FROM users WHERE email_hash = ?"
	var result DB_user
	err := m.db.Get(&result, query, emailhash)
	return result, err
}

func (m *mysqlEngine) GetUserByUsername(username string) (DB_user, error) {
	var query = "SELECT id, username, email_hash, email_crypt, password_hash, password_salt, bio, image FROM users WHERE username = ?"
	var result DB_user
	err := m.db.Get(&result, query, username)
	return result, err
}

func (m *mysqlEngine) GetUserById(id int) (DB_user, error) {
	var query = "SELECT id, username, email_hash, email_crypt, password_hash, password_salt, bio, image FROM users WHERE id = ?"
	var result DB_user
	err := m.db.Get(&result, query, id)
	return result, err
}

func (m *mysqlEngine) UpdateUser(id int) IUserUpdater {
	return &userUpdater{
		db:       m.db,
		id:       id,
		toUpdate: make(map[string]interface{}),
	}
}

/*
IUserUpdater
*/

type userUpdater struct {
	db       *sqlx.DB
	id       int
	toUpdate map[string]interface{}
}

func (u *userUpdater) Username(value string) IUserUpdater {
	u.toUpdate["username"] = value
	return u
}

func (u *userUpdater) Email_hash(value []byte) IUserUpdater {
	u.toUpdate["email_hash"] = value
	return u
}

func (u *userUpdater) Email_crypt(value []byte) IUserUpdater {
	u.toUpdate["email_crypt"] = value
	return u
}

func (u *userUpdater) Password_hash(value []byte) IUserUpdater {
	u.toUpdate["password_hash"] = value
	return u
}

func (u *userUpdater) Password_salt(value []byte) IUserUpdater {
	u.toUpdate["password_salt"] = value
	return u
}

func (u *userUpdater) Bio(value *string) IUserUpdater {
	u.toUpdate["bio"] = value
	return u
}

func (u *userUpdater) Image(value *string) IUserUpdater {
	u.toUpdate["image"] = value
	return u
}

func (u *userUpdater) Execute() error {
	if len(u.toUpdate) == 0 {
		return errors.New("Nothing to update!")
	}
	var where = make([]string, 0, len(u.toUpdate))
	var values = make([]interface{}, 0, len(u.toUpdate)+1)
	for key, value := range u.toUpdate {
		where = append(where, key+" = ?")
		values = append(values, value)
	}
	var query = "UPDATE users SET " + strings.Join(where, ", ") + " WHERE id = ?"
	values = append(values, u.id)
	_, err := u.db.Exec(query, values...)
	// 1062 duplicate entry?
	return err
}
