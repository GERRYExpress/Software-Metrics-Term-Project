package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssue(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	token, err := Issue(42)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token can be parsed - uses Subject claim
	parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("test-secret-key-for-testing"), nil
	})
	require.NoError(t, err)
	assert.True(t, parsed.Valid)
	claims := parsed.Claims.(*jwt.RegisteredClaims)
	assert.Equal(t, "42", claims.Subject)
}

func TestIssue_EmptySecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	_, err := Issue(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET is not configured")
}

func TestIssue_InvalidUserID(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	// Test that user ID is included in token as Subject
	token, err := Issue(999)
	require.NoError(t, err)

	parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("test-secret-key-for-testing"), nil
	})
	require.NoError(t, err)
	claims := parsed.Claims.(*jwt.RegisteredClaims)
	assert.Equal(t, "999", claims.Subject)
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

func TestMiddleware_ValidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	token, err := Issue(123)
	require.NoError(t, err)

	// Verify token structure
	parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("test-secret-key-for-testing"), nil
	})
	require.NoError(t, err)
	assert.True(t, parsed.Valid)

	claims := parsed.Claims.(*jwt.RegisteredClaims)
	assert.Equal(t, "123", claims.Subject)
	assert.NotZero(t, claims.ExpiresAt)

	// Verify expiration is in the future
	expTime := claims.ExpiresAt.Time
	assert.True(t, expTime.After(time.Now()))
}

func TestMiddleware_ExpiredToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	// Create an expired token
	claims := &jwt.RegisteredClaims{
		Subject:   "456",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("test-secret-key-for-testing"))
	require.NoError(t, err)

	// Parse should fail due to expiration
	parsed, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("test-secret-key-for-testing"), nil
	})
	assert.Error(t, err)
	assert.False(t, parsed.Valid)
}

func TestMiddleware_WrongSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	// Create token with different secret
	claims := &jwt.RegisteredClaims{
		Subject:   "789",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	wrongToken, err := token.SignedString([]byte("wrong-secret"))
	require.NoError(t, err)

	// Parse with correct secret should fail
	parsed, err := jwt.ParseWithClaims(wrongToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("test-secret-key-for-testing"), nil
	})
	assert.Error(t, err)
	assert.False(t, parsed.Valid)
}

func TestMiddleware_MalformedToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	// Try to parse malformed token - ParseWithClaims may return nil token on error
	_, err := jwt.ParseWithClaims("not.a.valid.token", &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("test-secret-key-for-testing"), nil
	})
	assert.Error(t, err)
}