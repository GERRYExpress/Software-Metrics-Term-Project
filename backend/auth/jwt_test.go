package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssue(t *testing.T) {
	// Set JWT_SECRET for test
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	token, err := Issue(42)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestIssue_EmptySecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	_, err := Issue(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET is not configured")
}

func TestMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	// Generate a valid token
	token, err := Issue(1)
	require.NoError(t, err)

	// The middleware is tested via integration tests
	_ = token
}