package shortener

import (
	"fmt"
	"net/http"
	"os"

	"github.com/asankov/shortener/internal/config"
	"github.com/asankov/shortener/internal/links"
	"github.com/asankov/shortener/internal/random"
	"github.com/asankov/shortener/internal/users"
	"golang.org/x/exp/slog"
)

type Shortener struct {
	server http.Server

	useSSL   bool
	certFile string
	keyFile  string

	logger *slog.Logger

	handler *handler

	configService ConfigService
}

type handler struct {
	db            Database
	userService   UserService
	authenticator Authenticator
	idGenerator   IDGenerator

	logger *slog.Logger
}

type Database interface {
	GetByID(id string) (*links.Link, error)
	GetAll() ([]*links.Link, error)

	Create(id string, url string) error
	Delete(id string) error
	IncrementClicks(id string) error
}

type IDGenerator interface {
	GenerateID() (string, error)
}

type UserService interface {
	GetUser(email, password string) (*users.User, error)
	CreateUser(email, password string, roles []users.Role) error
}

type Authenticator interface {
	NewTokenForUser(user *users.User) (string, error)
	DecodeToken(token string) (*users.User, error)
}

type ConfigService interface {
	ShouldCreateInitialUser() (bool, error)
}

func New(config *config.Config, db Database, idGenerator IDGenerator, userService UserService, authenticator Authenticator, configService ConfigService) (*Shortener, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s := &Shortener{
		useSSL:   config.UseSSL,
		certFile: config.SSL.CertFile,
		keyFile:  config.SSL.KeyFile,
		server: http.Server{
			Addr: fmt.Sprintf(":%d", config.Port),
		},
		logger: logger,
		handler: &handler{
			db:            db,
			userService:   userService,
			authenticator: authenticator,
			idGenerator:   idGenerator,
			logger:        logger,
		},
		configService: configService,
	}

	s.server.Handler = s.routes()

	return s, nil
}

func (s *Shortener) SetLogger(l *slog.Logger) *Shortener {
	s.logger = l
	s.handler.logger = l
	return s
}

func (s *Shortener) init() error {
	shouldCreateInitialUser, err := s.configService.ShouldCreateInitialUser()
	if err != nil {
		s.logger.Warn("error while checking whether to create initial user", "error", err)
	}
	if shouldCreateInitialUser {
		email, password := "admin@asankov.dev", random.Password(30)
		if err := s.handler.userService.CreateUser(email, password, []users.Role{users.RoleAdmin}); err != nil {
			return err
		}

		s.logger.Info(fmt.Sprintf("generated admin user with email [%s] and password [%s]", email, password))
	}

	return nil
}

func (s *Shortener) Start() error {
	if err := s.init(); err != nil {
		return err
	}

	if s.useSSL {
		s.logger.Info(fmt.Sprintf("Starting server on address [%s] with SSL\n", s.server.Addr))
		return s.server.ListenAndServeTLS(s.certFile, s.keyFile)
	}

	s.logger.Info(fmt.Sprintf("Starting server on address [%s] with no SSL\n", s.server.Addr))
	return s.server.ListenAndServe()
}
