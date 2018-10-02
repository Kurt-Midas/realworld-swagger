package session

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ISessionEngine interface {
	CreateUserToken(id int) (token string, err error)
	DecodeUserToken(token string) (id int, err error)
}

func NewJwtSessionEngine(secret []byte, seconds int) ISessionEngine {
	expiration := time.Duration(seconds) * time.Second
	return &jwtEngine{
		secret:     secret,
		expiration: expiration}
}

type rwjwt struct {
	jwt.StandardClaims
	rwclaims
}

type rwclaims struct {
	UserID int `json:"id"`
}

type jwtEngine struct {
	secret     []byte
	expiration time.Duration
}

func (eng *jwtEngine) CreateUserToken(id int) (string, error) {
	var claims = rwjwt{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(eng.expiration).Unix(),
		},
		rwclaims{
			UserID: id,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(eng.secret)
	if err != nil {
		return "", errors.New("Failed to sign JWT: " + err.Error())
	}
	return t, nil
}

func (eng *jwtEngine) DecodeUserToken(token string) (int, error) {
	parsed, err := jwt.ParseWithClaims(token, &rwjwt{}, func(tok *jwt.Token) (interface{}, error) {
		return eng.secret, nil
	})
	if err != nil {
		return 0, errors.New("Failed to decode signed JWT: " + err.Error())
	}

	if claims, ok := parsed.Claims.(*rwjwt); ok && parsed.Valid {
		return claims.UserID, nil
	} else {
		return 0, errors.New("Decoded JWT claims were invalid")
	}
}
