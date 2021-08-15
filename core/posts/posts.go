package posts

import (
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/service"
	"github.com/labstack/echo"
)

type Posts struct {
	ctx *service.Service
}

func (p *Posts) Setup(ctx *service.Service) {
	p.ctx = ctx

	postsRouter := p.ctx.Echo.Group("/api/v2/posts")

	// configure router

	// auth middleware
	// postsRouter.Use()

	// plural
	postsRouter.GET("", p.GetPostsHandler)
	postsRouter.POST("", p.CreatePostHandler)

	// singular
	postsRouter.GET("/:id", p.GetPostHandler)
	// postsRouter.PUT("/:id",)
	// postsRouter.DELETE("/:id",)

}

// Writes response based on accept header
// if header has application/json mime type as first index write response in xml
// else write response in json
func (p *Posts) ResponseWriter(c echo.Context, statusCode int, res interface{}) error {
	// check accept header
	accept := c.Request().Header["Accept"][0]

	// based on first value in accept header write response
	switch accept {
	case string(models.MimeTypesXML):
		// response with xml
		return c.XML(statusCode, res)
	default:
		// default response with json
		return c.JSON(statusCode, res)
	}
}
