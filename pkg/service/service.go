package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Tamplier2911/gorest/pkg/config"
	"github.com/Tamplier2911/gorest/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	Config *config.Config
	Logger *zap.SugaredLogger
	Server *http.Server
	Router *http.ServeMux

	// optional
	MySQL *gorm.DB
}

type InitializeOptions struct {
	MySQL bool
}

func (s *Service) Initialize(options *InitializeOptions) {
	var err error

	// create config
	s.Config = config.New()

	// create logger
	s.Logger = logger.
		New(s.Config.LogLevel, s.Config.Production).
		Named("Service")

	// create router
	s.Router = http.NewServeMux()

	// create server
	s.Server = &http.Server{
		Addr:           fmt.Sprintf(":%s", s.Config.Port),
		Handler:        s.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// create mysql connection with gorm package
	if options.MySQL {
		s.Logger.Infow("connecting to mysql server", "config", s.Config)
		s.MySQL, err = s.NewMySQL()
		if err != nil {
			s.Logger.Fatalw("failed to connect to mysql server", "config", s.Config, "err", err)
		}
		s.Logger.Infow("successfully connected to mysql server")
	}

}

func (s *Service) Start() {
	s.Logger.Infow(fmt.Sprintf("starting http server on port %s", s.Server.Addr))
	err := s.Server.ListenAndServe()
	if err != nil {
		s.Logger.Fatalw("failed to start server", "err", err)
	}
}
