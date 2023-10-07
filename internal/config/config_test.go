package config_test

import (
	"os"
	"testing"

	"github.com/asankov/shortener/internal/config"
	"github.com/stretchr/testify/require"
)

const secret = "abc123-secret"

func TestDefaults(t *testing.T) {
	setenv(t, "SHORTENER_SECRET", secret)

	config, err := config.NewFromEnv()

	require.NoError(t, err)
	require.Equal(t, 8080, config.Port)
	require.Equal(t, secret, config.Secret)
	require.False(t, config.ForceGenerateAdminUser)
}

func TestAllSet(t *testing.T) {
	setenv(t, "SHORTENER_PORT", "1234")
	setenv(t, "SHORTENER_SECRET", secret)
	setenv(t, "SHORTENER_FORCE_GENERATE_ADMIN_USER", "true")

	config, err := config.NewFromEnv()

	require.NoError(t, err)
	require.Equal(t, 1234, config.Port)
	require.Equal(t, secret, config.Secret)
	require.True(t, config.ForceGenerateAdminUser)
}

func TestRequired(t *testing.T) {
	_, err := config.NewFromEnv()

	require.Error(t, err)
}

func setenv(t *testing.T, key, value string) {
	t.Helper()

	err := os.Setenv(key, value)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := os.Unsetenv(key)
		require.NoError(t, err)
	})
}
