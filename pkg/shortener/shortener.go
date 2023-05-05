package shortener

import (
	"fmt"
	"net/http"

	"github.com/asankov/shortener/pkg/config"
	"github.com/asankov/shortener/pkg/links"
	"github.com/sirupsen/logrus"
)

type Shortener struct {
	server http.Server

	useSSL   bool
	certFile string
	keyFile  string

	logger *logrus.Logger

	db Database
}

type Database interface {
	GetByID(id string) (*links.Link, error)
	GetAll() ([]*links.Link, error)

	Create(id string, url string) error
	Delete(id string) error
}

func New(config *config.Config, db Database) (*Shortener, error) {
	s := &Shortener{
		useSSL:   config.UseSSL,
		certFile: config.SSL.CertFile,
		keyFile:  config.SSL.KeyFile,
		server: http.Server{
			Addr: fmt.Sprintf(":%d", config.Port),
		},
		logger: logrus.New(),
		db:     db,
	}

	s.server.Handler = s.routes()

	return s, nil
}

func (s *Shortener) SetLogger(l *logrus.Logger) *Shortener {
	s.logger = l
	return s
}

func (s *Shortener) Start() error {
	if s.useSSL {
		s.logger.Infof("Starting server on address [%s] with SSL\n", s.server.Addr)
		return s.server.ListenAndServeTLS(s.certFile, s.keyFile)
	}
	s.logger.Infof("Starting server on address [%s] with no SSL\n", s.server.Addr)
	return s.server.ListenAndServe()
}
