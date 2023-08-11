package main

import (
	"github.com/asankov/shortener/internal/dynamo"
	"github.com/asankov/shortener/pkg/config"
	"github.com/asankov/shortener/pkg/shortener"
	"golang.org/x/exp/slog"
)

func main() {
	if err := run(); err != nil {
		slog.Error("error while running shortener: %v", err)
	}
}

func run() error {

	config, err := config.NewFromEnv()
	if err != nil {
		return err
	}

	db, err := dynamo.New()
	if err != nil {
		return err
	}
	shortener, err := shortener.New(config, db, db)
	if err != nil {
		return err
	}

	return shortener.Start()
}
