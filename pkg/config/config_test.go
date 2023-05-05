package config_test

import (
	"os"
	"testing"

	"github.com/asankov/shortener/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestDefaults(t *testing.T) {
	config, err := config.NewFromEnv()

	require.NoError(t, err)
	require.Equal(t, false, config.UseSSL)
	require.Equal(t, 8080, config.Port)
	require.Equal(t, "", config.SSL.CertFile)
	require.Equal(t, "", config.SSL.KeyFile)
}

func TestAllSet(t *testing.T) {
	setenv(t, "SHORTENER_USE_SSL", "true")
	setenv(t, "SHORTENER_PORT", "1234")
	setenv(t, "SHORTENER_SSL_CERT_FILE", "cert.pem")
	setenv(t, "SHORTENER_SSL_KEY_FILE", "key.pem")

	config, err := config.NewFromEnv()

	require.NoError(t, err)
	require.Equal(t, true, config.UseSSL)
	require.Equal(t, 1234, config.Port)
	require.Equal(t, "cert.pem", config.SSL.CertFile)
	require.Equal(t, "key.pem", config.SSL.KeyFile)
}

func TestUseSSLSetButNoSSLConfig(t *testing.T) {
	setenv(t, "SHORTENER_USE_SSL", "true")

	_, err := config.NewFromEnv()

	require.Error(t, err)
	require.ErrorIs(t, err, config.ErrNoSSLConfig)
}

func setenv(t *testing.T, key, value string) {
	err := os.Setenv(key, value)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := os.Unsetenv(key)
		require.NoError(t, err)
	})
}
