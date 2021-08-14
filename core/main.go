package main

import (
	"github.com/Tamplier2911/gorest/core/comments"
	"github.com/Tamplier2911/gorest/core/posts"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/service"
)

type Monolith struct {
	service.Service
}

func (m *Monolith) Setup() {
	m.Initialize(&service.InitializeOptions{
		MySQL: true,
	})

	// default port '8080' || export GOREST_PORT='8080'
	// or manually
	m.Server.Addr = ":3000"

	// automigrate models
	m.Logger.Info("automigrating models")
	err := m.MySQL.AutoMigrate(&models.Post{}, &models.Comment{}, &models.User{})
	if err != nil {
		m.Logger.Fatalw("failed to automigrate models", "err", err)
	}

	// v1 posts
	posts := posts.Posts{}
	posts.Setup(&m.Service)

	// v1 comments
	comments := comments.Comments{}
	comments.Setup(&m.Service)

}

func main() {
	s := Monolith{}
	s.Setup()
	s.Start()
}
