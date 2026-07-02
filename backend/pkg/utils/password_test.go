package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHashRoundTrip(t *testing.T) {
	hash, err := HashPassword("S3cret-passw0rd!")
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	assert.NotEqual(t, "S3cret-passw0rd!", hash, "hash must not be the plaintext")

	assert.True(t, CheckPassword("S3cret-passw0rd!", hash), "correct password should verify")
	assert.False(t, CheckPassword("wrong-password", hash), "wrong password must not verify")
}

func TestCheckPasswordRejectsInvalidHash(t *testing.T) {
	assert.False(t, CheckPassword("anything", "not-a-bcrypt-hash"))
	assert.False(t, CheckPassword("anything", ""))
}
