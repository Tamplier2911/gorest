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

	// TODO: consider refactoring this and abstracting into a router
	// configure router
	s.Router.HandleFunc("/v1/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.GetPostsHandler(w, r)
		case http.MethodPost:
			s.CreatePostHandler(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	s.Router.HandleFunc("/v1/posts/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.GetPostHandler(w, r)
		case http.MethodPut:
			s.UpdatePostHandler(w, r)
		case http.MethodDelete:
			s.DeletePostHandler(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

}

func main() {
	s := Monolith{}
	s.Setup()
	s.Start()
}
