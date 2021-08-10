package main

import (
	"net/http"

	"github.com/Tamplier2911/gorest/core/posts"
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
	err := s.MySQL.AutoMigrate(&posts.Post{})
	if err != nil {
		s.Logger.Fatalw("failed to automigrate models", "err", err)
	}

	// posts
	posts := posts.Posts{}
	posts.Setup(&s.Service)

	// configure router
	s.Router.HandleFunc("/v1/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			posts.GetPostsHandler(w, r)
		case http.MethodPost:
			posts.CreatePostHandler(w, r)
		}
	})

	s.Router.HandleFunc("/v1/posts/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			posts.GetPostHandler(w, r)
		case http.MethodPut:
			posts.UpdatePostHandler(w, r)
		case http.MethodDelete:
			posts.DeletePostHandler(w, r)
		}
	})

	// comments

	// TODO: create cruds for comments

}

func main() {
	s := Monolith{}
	s.Setup()
	s.Start()
}
