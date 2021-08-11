package main

import (
	"net/http"

	"github.com/Tamplier2911/gorest/core/comments"
	"github.com/Tamplier2911/gorest/core/posts"
	"github.com/Tamplier2911/gorest/pkg/service"
)

type Monolith struct {
	service.Service
}

func (m *Monolith) Setup() {
	m.Initialize(&service.InitializeOptions{
		MySQL: true,
	})

	// set port - default 8080 || export GOREST_PORT='3000'
	m.Server.Addr = ":3000"

	// automigrate models
	m.Logger.Info("automigrating models")
	// err := m.MySQL.AutoMigrate(&models.{})
	err := m.MySQL.AutoMigrate(&posts.Post{}, &comments.Comment{})
	if err != nil {
		m.Logger.Fatalw("failed to automigrate models", "err", err)
	}

	// posts
	posts := posts.Posts{}
	posts.Setup(&m.Service)

	// configure router
	m.Router.HandleFunc("/v1/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			posts.GetPostsHandler(w, r)
		case http.MethodPost:
			posts.CreatePostHandler(w, r)
		}
		w.WriteHeader(http.StatusNotFound)
	})

	m.Router.HandleFunc("/v1/posts/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			posts.GetPostHandler(w, r)
		case http.MethodPut:
			posts.UpdatePostHandler(w, r)
		case http.MethodDelete:
			posts.DeletePostHandler(w, r)
		}
		w.WriteHeader(http.StatusNotFound)
	})

	// comments
	comments := comments.Comments{}
	comments.Setup(&m.Service)

	// configure router
	m.Router.HandleFunc("/v1/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// comments.GetPostsHandler(w, r)
		case http.MethodPost:
			comments.CreateCommentHandler(w, r)
		}
		w.WriteHeader(http.StatusNotFound)
	})

	m.Router.HandleFunc("/v1/comments/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// comments.GetPostHandler(w, r)
		case http.MethodPut:
			// comments.UpdatePostHandler(w, r)
		case http.MethodDelete:
			// comments.DeletePostHandler(w, r)
		}
		w.WriteHeader(http.StatusNotFound)
	})

}

func main() {
	s := Monolith{}
	s.Setup()
	s.Start()
}
