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

	// set port - default 8080 || export GOREST_PORT='3000'
	s.Server.Addr = ":3000"

	// automigrate models
	s.Logger.Info("automigrating models")
	err := s.MySQL.AutoMigrate(&Post{})
	if err != nil {
		s.Logger.Fatalw("failed to automigrate models", "err", err)
	}

	// TODO: is it possible to add all posts endpoints to sub directory to abstract it away
	// TODO: consider refactoring this and abstracting into a router
	// TODO: create cruds for comments

	// configure router
	s.Router.HandleFunc("/v1/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.GetPostsHandler(w, r)
		case http.MethodPost:
			s.CreatePostHandler(w, r)
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
		}
	})

}

func main() {
	s := Monolith{}
	s.Setup()
	s.Start()
}
