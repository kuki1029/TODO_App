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
func Verify(hash string, raw string) (bool, error) {
	isSame, err := argon2id.ComparePasswordAndHash(hash, raw)
	return isSame, err

}
