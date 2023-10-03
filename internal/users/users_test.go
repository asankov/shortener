package users_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/asankov/shortener/internal/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsers(t *testing.T) {
	user := users.User{
		Email: "asankov@asankov.dev",
		Roles: []users.Role{users.RoleUser},
	}

	require.True(t, user.HasRole(users.RoleUser))
	require.False(t, user.HasRole(users.RoleAdmin))
}

func TestRoleFrom(t *testing.T) {
	testCases := []struct {
		name     string
		variants []string
		role     users.Role
	}{
		{
			name:     "ADMIN",
			variants: []string{"admin", "ADMIN", "Admin", "aDmin", "adMIN"},
			role:     users.RoleAdmin,
		},
		{
			name:     "USER",
			variants: []string{"user", "USER", "User", "uSer", "usER"},
			role:     users.RoleUser,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, variant := range testCase.variants {
				role, err := users.RoleFrom(variant)
				if assert.NoError(t, err) {
					assert.Equal(t, testCase.role, role)
				}
			}
		})
	}

	t.Run("TestInvalidRole", func(t *testing.T) {
		_, err := users.RoleFrom("unknown")

		require.Error(t, err)
		require.True(t, errors.Is(err, users.ErrInvalidRole))
	})
}

func TestRole(t *testing.T) {
	for _, role := range []users.Role{users.RoleAdmin, users.RoleUser, users.Role(5)} {
		res, err := json.Marshal(role)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(`"%s"`, role.String()), string(res))
	}
}
