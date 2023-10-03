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
	require.Equal(t, false, config.UseSSL)
	require.Equal(t, 8080, config.Port)
	require.Equal(t, "", config.SSL.CertFile)
	require.Equal(t, "", config.SSL.KeyFile)
	require.Equal(t, secret, config.Secret)
}

func TestAllSet(t *testing.T) {
	setenv(t, "SHORTENER_USE_SSL", "true")
	setenv(t, "SHORTENER_PORT", "1234")
	setenv(t, "SHORTENER_SSL_CERT_FILE", "cert.pem")
	setenv(t, "SHORTENER_SSL_KEY_FILE", "key.pem")
	setenv(t, "SHORTENER_SECRET", secret)

	config, err := config.NewFromEnv()

	require.NoError(t, err)
	require.Equal(t, true, config.UseSSL)
	require.Equal(t, 1234, config.Port)
	require.Equal(t, "cert.pem", config.SSL.CertFile)
	require.Equal(t, "key.pem", config.SSL.KeyFile)
	require.Equal(t, secret, config.Secret)
}

func TestUseSSLSetButNoSSLConfig(t *testing.T) {
	setenv(t, "SHORTENER_USE_SSL", "true")
	setenv(t, "SHORTENER_SECRET", secret)

	_, err := config.NewFromEnv()

	require.Error(t, err)
	require.ErrorIs(t, err, config.ErrNoSSLConfig)
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
