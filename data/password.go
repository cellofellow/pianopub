package data

import (
	"crypto/rand"
	"crypto/sha256"

	"code.google.com/p/go.crypto/pbkdf2"
)

type HashedPassword struct {
	Salt string
	Hash string
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
	return HashedPassword{Salt: string(salt), Hash: string(hash)}
}

func CheckPassword(password, salt string) HashedPassword {
	hash := pbkdf2.Key([]byte(password), []byte(salt), 8192, 32, sha256.New)
	return HashedPassword{Salt: salt, Hash: string(hash)}
}
