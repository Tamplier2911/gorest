package main

import (
	"github.com/Tamplier2911/gorest/core/comments"
	"github.com/Tamplier2911/gorest/core/posts"
	v1_comments "github.com/Tamplier2911/gorest/core/v1_comments"
	v1_posts "github.com/Tamplier2911/gorest/core/v1_posts"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/service"

	_ "github.com/Tamplier2911/gorest/core/docs"
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

	// /api/v1/posts
	deprecatedPosts := v1_posts.Posts{}
	deprecatedPosts.Setup(&m.Service)
	// /api/v1/comments
	deprecatedComments := v1_comments.Comments{}
	deprecatedComments.Setup(&m.Service)

	// m.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	// /api/v2/posts
	posts := posts.Posts{}
	posts.Setup(&m.Service)
	// /api/v2/comments
	comments := comments.Comments{}
	comments.Setup(&m.Service)
}

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @title Go REST API example
// @version 2.0
// @description This is a sample rest api realized in go language for education purposes.
//
// @contact.email artyom.nikolaev@syahoo.com
//
// @host localhost:8000
// @BasePath /api/v2
// @query.collection.format multi
func main() {
	m := Monolith{}
	m.Setup()
	m.Start()
}
