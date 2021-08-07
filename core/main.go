package main

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/service"
)

type Monolith struct {
	service.Service
}

func (s *Monolith) Setup() {
	s.Initialize(&service.InitializeOptions{
		MySQL: true,
	})

	// set port - default 8080
	s.Server.Addr = ":3000"

	// automigrate models
	s.Logger.Info("automigrating models")
	err := s.MySQL.AutoMigrate(&Post{})
	if err != nil {
		s.Logger.Fatalw("failed to automigrate models", "err", err)
	}

	// configure router
	s.Router.HandleFunc("/v1/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// handle get one
			// w.WriteHeader(http.StatusOK)
			// handle get all with limit and offset
			w.WriteHeader(http.StatusOK)
		case http.MethodPost:
			// handle creat one
			s.CreatePostHandler(w, r)
		case http.MethodPut:
			// handle update one
			w.WriteHeader(http.StatusOK)
		case http.MethodDelete:
			// handle delete one
			w.WriteHeader(http.StatusNoContent)
		default:
			// bad request
			w.WriteHeader(http.StatusBadRequest)
		}
	})

}

func main() {
	s := Monolith{}
	s.Setup()
	s.Start()
}
