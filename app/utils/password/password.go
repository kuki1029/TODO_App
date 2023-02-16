package password

import (
	"github.com/alexedwards/argon2id"
)

// Generate return a hashed password
func Generate(raw string) string {
	hash, err := argon2id.CreateHash(raw, argon2id.DefaultParams)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

// Verify compares a hashed password with plaintext password
func Verify(raw string, hash string) (bool, error) {
	isSame, err := argon2id.ComparePasswordAndHash(raw, hash)
	return isSame, err

}
