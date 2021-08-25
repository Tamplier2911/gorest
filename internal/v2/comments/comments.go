package comments

import (
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/service"
	"github.com/labstack/echo/v4"
)

type Comments struct {
	*service.Service
}

func (cm *Comments) Setup(s *service.Service) {
	cm.Service = s

	// configure router
	CommentsRouter := cm.Echo.Group("/api/v2/comments")

	CommentsRouter.GET("", cm.GetCommentsHandler)
	CommentsRouter.POST("", service.AuthenticationMiddleware(cm.Logger, cm.Config, cm.CreateCommentHandler))
	CommentsRouter.GET("/:id", cm.GetCommentHandler)
	CommentsRouter.PUT("/:id", service.AuthenticationMiddleware(cm.Logger, cm.Config, cm.UpdateCommentHandler))
	CommentsRouter.DELETE("/:id", service.AuthenticationMiddleware(cm.Logger, cm.Config, cm.DeleteCommentHandler))
}

// Writes response based on accept header
// if header has application/xml mime type as first index, write response in xml else write response in json
func (p *Comments) ResponseWriter(c echo.Context, statusCode int, res interface{}) error {
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
