package internal

import (
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/service"

	_ "github.com/Tamplier2911/gorest/internal/docs"
	v1comments "github.com/Tamplier2911/gorest/internal/v1/comments"
	v1posts "github.com/Tamplier2911/gorest/internal/v1/posts"
	"github.com/Tamplier2911/gorest/internal/v2/auth"
	"github.com/Tamplier2911/gorest/internal/v2/comments"
	"github.com/Tamplier2911/gorest/internal/v2/posts"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Monolith struct {
	service.Service
}

func (m *Monolith) Setup() {
	m.Initialize(&service.InitializeOptions{
		MySQL: true,
		Echo:  true,
	})

	// default port '8080' || export GOREST_PORT='8080' || m.Server.Addr = ":3000"

	// automigrate models
	m.Logger.Info("automigrating models")
	err := m.MySQL.AutoMigrate(
		&models.Post{},
		&models.Comment{},
		&models.User{},
		&models.AuthProvider{},
	)
	if err != nil {
		m.Logger.Fatalw("failed to automigrate models", "err", err)
	}

	// /swagger/index.html
	m.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	// /api/v1/posts
	v1posts.Posts{}.Setup(&m.Service)
	// /api/v1/comments
	v1comments.Comments{}.Setup(&m.Service)

	// /api/v2/auth
	auth.Auth{}.Setup(&m.Service)
	// /api/v2/posts
	posts.Posts{}.Setup(&m.Service)
	// /api/v2/comments
	comments.Comments{}.Setup(&m.Service)
}
