package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Tamplier2911/gorest/pkg"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	Config *config.Config
	Logger *zap.SugaredLogger
	Server *http.Server
	// Router *gin.RouterGroup

	// optional
	MySQL *gorm.DB
}

type InitializeOptions struct {
	MySQL bool
}

func (s *Service) Initialize(options *InitializeOptions) {
	// read config
	s.Config = config.New()

	// create logger
	s.Logger = logger.
		New(s.Config.LogLevel, s.Config.Production).
		Named("Service")

	// create server
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", s.Config.Port),
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.Server = server

	// create router
	// s.Router = 

	// create mysql connection with gorm package
	if options.MySQL {
		s.Logger.Infow("connecting to postgresql", "config", s.Config)
		db, err := s.NewMySQL()
		if err != nil {
			s.Logger.Fatalw("failed to connect to postgresql", "config", s.Config, "err", err)
		}
		s.MySQL = db
		s.Logger.Debugw("connected to postgresql")
	}

}

func (s *Service) Start() {
	err := s.Server.ListenAndServe()
	if err != nil {
		s.Logger.Fatalw("failed to start server", "err", err)
	}
}
