package main

import (
	"github.com/Tamplier2911/gorest/core/posts"
	v1_comments "github.com/Tamplier2911/gorest/core/v1_comments"
	v1_posts "github.com/Tamplier2911/gorest/core/v1_posts"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/service"
)

type Monolith struct {
	service.Service
}

func (m *Monolith) Setup() {
	m.Initialize(&service.InitializeOptions{
		MySQL: true,
		Echo:  true,
	})

	// default port '8080' || export GOREST_PORT='8080'
	// or manually
	// m.Server.Addr = ":3000"

	// automigrate models
	m.Logger.Info("automigrating models")
	err := m.MySQL.AutoMigrate(&models.Post{}, &models.Comment{}, &models.User{})
	if err != nil {
		m.Logger.Fatalw("failed to automigrate models", "err", err)
	}

	// v1
	// posts
	deprecatedPosts := v1_posts.Posts{}
	deprecatedPosts.Setup(&m.Service)
	// comments
	deprecatedComments := v1_comments.Comments{}
	deprecatedComments.Setup(&m.Service)

	//v2
	// posts
	posts := posts.Posts{}
	posts.Setup(&m.Service)
}

func main() {
	s := Monolith{}
	s.Setup()
	s.Start()
}
