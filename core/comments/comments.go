package comments

import (
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/service"
	"github.com/labstack/echo/v4"
)

type Comments struct {
	*service.Service
}

// @BasePath /api/v2/comments
func (cm *Comments) Setup(service *service.Service) {
	cm.Service = service

	// configure router
	CommentsRouter := cm.Echo.Group("/api/v2/comments")

	// auth middleware
	// TODO: only owners can remove and update comments
	// CommentsRouter.Use()

	CommentsRouter.GET("", cm.GetCommentsHandler)
	CommentsRouter.POST("", cm.CreateCommentHandler)
	CommentsRouter.GET("/:id", cm.GetCommentHandler)
	CommentsRouter.PUT("/:id", cm.UpdateCommentHandler)
	CommentsRouter.DELETE("/:id", cm.DeleteCommentHandler)
}

// Writes response based on accept header
// if header has application/xml mime type as first index, write response in xml else write response in json
func (p *Comments) ResponseWriter(c echo.Context, statusCode int, res interface{}) error {
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
