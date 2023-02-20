package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests the verify function for passwords. We do not test the Generate function
// as that mainly calls the argon2id functions which are already tested.
// Here we test both cases where the passwords match and a case where they
// don't match
func TestVerify(t *testing.T) {
	pass := "password"
	pass2 := "wrongpassword"
	hash := Generate(pass)
	isSame, err := Verify(pass, hash)
	isSame2, err2 := Verify(pass2, hash)
	assert.Nilf(t, err, "Expected password to match hash")
	// Only returns error if problem with verifying.
	assert.Nilf(t, err2, "Expected password to not match hash")
	assert.True(t, isSame)
	assert.False(t, isSame2)
}
