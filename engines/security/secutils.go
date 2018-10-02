package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

type ISecurityEngine interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	Hash(value []byte, salt []byte) [32]byte
	NewSalt() ([]byte, error)
}

type securityEngine struct {
	aeskey      [32]byte
	defaultSalt []byte
}

func NewSecurityEngine(aeskey [32]byte, defaultSalt []byte) ISecurityEngine {
	return &securityEngine{aeskey, defaultSalt}
}

func (engine *securityEngine) Encrypt(plaintext []byte) ([]byte, error) {
	block, blockErr := aes.NewCipher(engine.aeskey[:])
	if blockErr != nil {
		return nil, errors.New(blockErr.Error())
	}
	// Nonces should be 96 bits
	nonce := make([]byte, 12)
	if _, randErr := io.ReadFull(rand.Reader, nonce); randErr != nil {
		return nil, errors.New(randErr.Error())
	}

	aesgcm, gcmErr := cipher.NewGCM(block)
	if gcmErr != nil {
		return nil, errors.New(gcmErr.Error())
	}
	//dst, nonce, plaintext, additionalData
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}

func (engine *securityEngine) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < 12 {
		return nil, errors.New("Invalid ciphertext")
	}
	block, blockErr := aes.NewCipher(engine.aeskey[:])
	if blockErr != nil {
		return nil, errors.New(blockErr.Error())
	}
	aesgcm, gcmErr := cipher.NewGCM(block)
	if gcmErr != nil {
		return nil, errors.New(gcmErr.Error())
	}
	nonce := ciphertext[:12]     //first 12 are nonce
	ciphertext = ciphertext[12:] //last len-12 are ciphertext
	plaintext, openErr := aesgcm.Open(nil, nonce, ciphertext, nil)
	if openErr != nil {
		return nil, errors.New(openErr.Error())
	}
	return plaintext, nil
}

func (engine *securityEngine) Hash(value []byte, salt []byte) [32]byte {
	if salt == nil {
		salt = engine.defaultSalt
	}
	return sha256.Sum256(append(value, salt...))
}

func (engine *securityEngine) NewSalt() ([]byte, error) {
	salt := make([]byte, 32)
	if _, randErr := io.ReadFull(rand.Reader, salt); randErr != nil {
		return nil, errors.New(randErr.Error())
	}
	return salt, nil
}
