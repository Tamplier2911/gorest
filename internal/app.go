package app

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

type Application struct {
	service.Service
}

func (a *Application) Setup() {
	a.Initialize(&service.InitializeOptions{
		MySQL:     true,
		Echo:      true,
		Validator: true,
	})

	// default port '8080' || export GOREST_PORT='8080' || m.Server.Addr = ":3000"

	// automigrate models
	a.Logger.Info("automigrating models")
	err := a.MySQL.AutoMigrate(
		&models.User{},
		&models.AuthProvider{},
		&models.Post{},
		&models.Comment{},
	)
	if err != nil {
		a.Logger.Fatalw("failed to automigrate models", "err", err)
	}

	// /swagger/index.html
	a.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	// /api/v1/posts
	v1posts.Posts{}.Setup(&a.Service)
	// /api/v1/comments
	v1comments.Comments{}.Setup(&a.Service)

	// /api/v2/auth
	auth.Auth{}.Setup(&a.Service)
	// /api/v2/posts
	posts.Posts{}.Setup(&a.Service)
	// /api/v2/comments
	comments.Comments{}.Setup(&a.Service)
}
