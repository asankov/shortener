package config

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
)

var ErrNoSSLConfig = errors.New("no SSL config provided")

type Config struct {
	UseSSL bool `default:"false" split_words:"true"`
	Port   int  `default:"8080"`
	SSL    SSL
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

	if config.UseSSL && (config.SSL.CertFile == "" || config.SSL.KeyFile == "") {
		return nil, ErrNoSSLConfig
	}

	return &config, nil
}
