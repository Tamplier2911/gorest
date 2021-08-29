package service

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Tamplier2911/gorest/pkg/config"
	"github.com/Tamplier2911/gorest/pkg/logger"
	"github.com/labstack/echo/v4"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	Config *config.Config
	Logger *zap.SugaredLogger

	// default server and multiplexer
	Server *http.Server
	Router *http.ServeMux

	// optional
	MySQL     *gorm.DB
	Echo      *echo.Echo
	Validator *validator.Validate
}

type InitializeOptions struct {
	MySQL     bool
	Echo      bool
	Validator bool
}

func (s *Service) Initialize(options *InitializeOptions) {
	var err error

	// create config
	s.Config = config.New()

	// create logger
	s.Logger = logger.
		New(s.Config.LogLevel, s.Config.Production).
		Named("Service")

	// get port
	port := fmt.Sprintf(":%s", s.Config.Port)
	if s.Config.Production {
		port = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	// create default router
	s.Router = http.NewServeMux()

	// create default server
	s.Server = &http.Server{
		Addr:           port,
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

	// create echo instance
	if options.Echo {
		s.Logger.Infow("wiring echo framework server")
		s.Echo = echo.New()
	}

	// create validator
	if options.Validator {
		s.Logger.Infow("wiring validator")
		s.Validator = validator.New()
	}
}

func (s *Service) Start() {
	// if echo initialized run both servers in parallel
	if s.Echo != nil {
		var wg sync.WaitGroup

		// create error channels
		defaultServerError := make(chan error, 1)
		echoServerError := make(chan error, 1)

		// start default http server
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			s.Logger.Infow(fmt.Sprintf("starting default http server - base url: %s port: %s", s.Config.BaseURL, s.Server.Addr))
			defaultServerError <- s.Server.ListenAndServe()
			// cleanup
		}(&wg)

		// start echo server
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			// TODO: add dynamic port
			port := "8000"
			s.Logger.Infow(fmt.Sprintf("starting echo http server - base url: %s port: %s", s.Config.BaseURL, port))
			echoServerError <- s.Echo.Start(fmt.Sprintf(":%s", port))
			// cleanup
		}(&wg)

		select {
		case err := <-defaultServerError:
			// handle error and close echo server
			s.Echo.Close()
			s.Logger.Fatalw("default server error:", "err", err)
		case err := <-echoServerError:
			// handle error and close default server
			s.Server.Close()
			s.Logger.Fatalw("echo server error:", "err", err)
		}

		// cleanup

		wg.Wait()
	} else {
		// else just run default http server
		s.Logger.Infow(fmt.Sprintf("starting default http server - base url: %s port: %s", s.Config.BaseURL, s.Server.Addr))
		err := s.Server.ListenAndServe()
		if err != nil {
			s.Logger.Fatalw("failed to start server", "err", err)
		}
	}
}
