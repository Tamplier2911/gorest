package posts

import (
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/service"
	"github.com/labstack/echo/v4"
)

type Posts struct {
	*service.Service
}

func (p Posts) Setup(s *service.Service) {
	p.Service = s

	// configure router
	PostsRouter := p.Echo.Group("/api/v2/posts")

	PostsRouter.GET("", p.GetPostsHandler)

	PostsRouter.POST("", service.AuthenticationMiddleware(p.Logger, p.Config, p.CreatePostHandler))
	PostsRouter.GET("/:id", p.GetPostHandler)
	PostsRouter.PUT("/:id", service.AuthenticationMiddleware(p.Logger, p.Config, p.UpdatePostHandler))
	PostsRouter.DELETE("/:id", service.AuthenticationMiddleware(p.Logger, p.Config, p.DeletePostHandler))
}

// Writes response based on accept header
// if header has application/xml mime type as first index, write response in xml else write response in json
func (p *Posts) ResponseWriter(c echo.Context, statusCode int, res interface{}) error {
	// check accept header
	accept := c.Request().Header["Accept"]
	if len(accept) == 0 {
		// default response if accept header is not provided
		return c.JSON(statusCode, res)
	}

	// based on first value in accept header write response
	switch accept[0] {
	case string(models.MimeTypesXML):
		// response with xml
		return c.XML(statusCode, res)
	default:
		// default response with json
		return c.JSON(statusCode, res)
	}
}
