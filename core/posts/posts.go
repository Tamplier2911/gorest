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

	// configure router
	PostsRouter := p.ctx.Echo.Group("/api/v2/posts")

	// auth middleware
	// TODO: only owners can remove and update posts
	// PostsRouter.Use()

	PostsRouter.GET("", p.GetPostsHandler)
	PostsRouter.POST("", p.CreatePostHandler)
	PostsRouter.GET("/:id", p.GetPostHandler)
	PostsRouter.PUT("/:id", p.UpdatePostHandler)
	PostsRouter.DELETE("/:id", p.DeletePostHandler)
}

// Writes response based on accept header
// if header has application/xml mime type as first index, write response in xml else write response in json
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
