package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"code.google.com/p/go.crypto/pbkdf2"
)

type HashedPassword struct {
	Salt string `json:"-"`
	Hash string `json:"-"`
}

func HashPassword(password string) HashedPassword {
	var err error
	salt := make([]byte, 12)
	_, err = rand.Read(salt)
	if err != nil {
		// Reading from /dev/urandom really should not error.
		panic(err.Error())
	}
	hash := pbkdf2.Key([]byte(password), salt, 8192, 32, sha256.New)
	return HashedPassword{
		Salt: base64.StdEncoding.EncodeToString(salt),
		Hash: base64.StdEncoding.EncodeToString(hash),
	}
}

func CheckPassword(password, salt string) HashedPassword {
	s, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return HashedPassword{}
	}
	hash := pbkdf2.Key([]byte(password), s, 8192, 32, sha256.New)
	return HashedPassword{
		Salt: salt,
		Hash: base64.StdEncoding.EncodeToString(hash),
	}
}
