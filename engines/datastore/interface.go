package datastore

type IDatastoreEngine interface {
	NewUser(username string, email_hash, email_cipher, pwhash, salt []byte) error
	GetUserByEmailHash([]byte) (DB_user, error)
	GetUserByUsername(string) (DB_user, error)
	GetUserById(int) (DB_user, error)
	UpdateUser(id int) IUserUpdater
}

type DB_user struct {
	Id            int     `db:"id"`
	Username      string  `db:"username"`
	Email_hash    []byte  `db:"email_hash"`
	Email_crypt   []byte  `db:"email_crypt"`
	Password_hash []byte  `db:"password_hash"`
	Password_salt []byte  `db:"password_salt"`
	Bio           *string `db:"bio"`
	Image         *string `db:"image"`
}

type IUserUpdater interface {
	Username(value string) IUserUpdater
	Email_hash(value []byte) IUserUpdater
	Email_crypt(value []byte) IUserUpdater
	Password_hash(value []byte) IUserUpdater
	Password_salt(value []byte) IUserUpdater
	Bio(value *string) IUserUpdater
	Image(value *string) IUserUpdater
	Execute() error
}
