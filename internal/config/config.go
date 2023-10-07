package config

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
)

var ErrNoSSLConfig = errors.New("no SSL config provided")

type Config struct {
	Port int `default:"8080"`

	// Secret is the secret used to generate the JWT token.
	Secret string `required:"true"`

	UseInMemoryDB bool `envconfig:"SHORTENER_USE_IN_MEMORY_DB"`

	ForceGenerateAdminUser bool `split_words:"true"`
}

type SSL struct {
	CertFile string `split_words:"true"`
	KeyFile  string `split_words:"true"`
}

func NewFromEnv() (*Config, error) {
	var config Config
	if err := envconfig.Process("SHORTENER", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
