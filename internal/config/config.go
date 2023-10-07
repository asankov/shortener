package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config is the struct that configures the service.
type Config struct {
	// Port controls on which port the service will listen to.
	Port int `default:"8080"`
	// Secret is the secret used to generate the JWT token.
	Secret string `required:"true"`
	// UseInMemoryDB controls whether an in-memory DB will be used for the service.
	//
	// This is useful for local testing, but not for production use.
	UseInMemoryDB bool `envconfig:"SHORTENER_USE_IN_MEMORY_DB"`
	// ForceGenerateAdminUser controls whether or not to ALWAYS generate an admin user on startup.
	//
	// If true, an admin user will be created on startup.
	// If false, an admin user might be created due to other conditions.
	ForceGenerateAdminUser bool `split_words:"true"`
}

// NewFromEnv creates new config with values loaded from environment variables.
func NewFromEnv() (*Config, error) {
	var config Config
	if err := envconfig.Process("SHORTENER", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
