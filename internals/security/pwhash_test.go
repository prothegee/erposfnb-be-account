package internal_security_test

import (
	"testing"

	isPW "github.com/prothegee/erposfnb-be-account/internals/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Password_Hasher(t *testing.T) {
	raw := "superPassword123!@"
	wrongPassword := "notMyPassword123!@"

	// 1. Hash the real password
	hash, err := isPW.PasswordHashArgon2(raw)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	// Sanity: hash must NOT equal the raw password
	assert.NotEqual(t, raw, hash)

	// 2. Correct password verifies OK
	ok, err := isPW.PasswordVerifyArgon2(raw, hash)
	require.NoError(t, err)
	assert.True(t, ok, "correct password should verify")

	// 3. Wrong password fails to verify
	ok, err = isPW.PasswordVerifyArgon2(wrongPassword, hash)
	require.NoError(t, err)
	assert.False(t, ok, "wrong password should NOT verify")
}
