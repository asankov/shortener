package shortener

import (
	"fmt"
	"net/http"
	"os"

	"github.com/asankov/shortener/internal/users"
	"github.com/asankov/shortener/pkg/config"
	"github.com/asankov/shortener/pkg/links"
	"golang.org/x/exp/slog"
)

type Shortener struct {
	server http.Server

	useSSL   bool
	certFile string
	keyFile  string

	logger *slog.Logger

	db            Database
	idGenerator   IDGenerator
	userService   UserService
	authenticator Authenticator
}

type Database interface {
	GetByID(id string) (*links.Link, error)
	GetAll() ([]*links.Link, error)

	Create(id string, url string) error
	Delete(id string) error
}

type IDGenerator interface {
	GenerateID() (string, error)
}

type UserService interface {
	Get(email, password string) (*users.User, error)
}

type Authenticator interface {
	NewTokenForUser(user *users.User) (string, error)
}

func New(config *config.Config, db Database, idGenerator IDGenerator, userService UserService, authenticator Authenticator) (*Shortener, error) {
	s := &Shortener{
		useSSL:   config.UseSSL,
		certFile: config.SSL.CertFile,
		keyFile:  config.SSL.KeyFile,
		server: http.Server{
			Addr: fmt.Sprintf(":%d", config.Port),
		},
		logger:        slog.New(slog.NewTextHandler(os.Stdout, nil)),
		db:            db,
		idGenerator:   idGenerator,
		userService:   userService,
		authenticator: authenticator,
	}

	s.server.Handler = s.routes()

	return s, nil
}

func (s *Shortener) SetLogger(l *slog.Logger) *Shortener {
	s.logger = l
	return s
}

func (s *Shortener) Start() error {
	if s.useSSL {
		s.logger.Info(fmt.Sprintf("Starting server on address [%s] with SSL\n", s.server.Addr))
		return s.server.ListenAndServeTLS(s.certFile, s.keyFile)
	}

	s.logger.Info(fmt.Sprintf("Starting server on address [%s] with no SSL\n", s.server.Addr))
	return s.server.ListenAndServe()
}
