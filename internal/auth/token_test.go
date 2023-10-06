package auth_test

import (
	"testing"
	"time"

	"github.com/asankov/shortener/internal/auth"
	"github.com/asankov/shortener/internal/users"
	"github.com/stretchr/testify/require"
)

const secret = "ABC123"

var user = &users.User{
	Email: "admin@asankov.dev",
	Roles: []users.Role{users.RoleAdmin},
}

func TestAuthenticator(t *testing.T) {
	authenticator := auth.NewAutheniticator(secret)

	token, err := authenticator.NewTokenForUser(user)
	require.NoError(t, err)

	decodedUser, err := authenticator.DecodeToken(token)
	require.NoError(t, err)
	require.Equal(t, user.Email, decodedUser.Email)
	require.Len(t, decodedUser.Roles, 1)
	require.Equal(t, users.RoleAdmin, decodedUser.Roles[0])

	t.Run("TestExpiredToken", func(t *testing.T) {
		token, err := authenticator.NewTokenForUserWithExpiration(user, -5*time.Minute)
		require.NoError(t, err)

		_, err = authenticator.DecodeToken(token)
		require.Error(t, err)
		require.ErrorIs(t, err, auth.ErrTokenExpired)
	})

	t.Run("TestInvalidSignature", func(t *testing.T) {
		token, err := auth.NewAutheniticator("another-secret").NewTokenForUser(user)
		require.NoError(t, err)

		_, err = authenticator.DecodeToken(token)
		require.Error(t, err)
		require.ErrorIs(t, err, auth.ErrInvalidSignature)
	})

	t.Run("TestInvalidFormat", func(t *testing.T) {
		_, err := authenticator.DecodeToken("abc.xyz")

		require.Error(t, err)
		require.ErrorIs(t, err, auth.ErrInvalidFormat)
	})
}
