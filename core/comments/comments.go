package comments

import (
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/Tamplier2911/gorest/pkg/service"
	"github.com/labstack/echo"
)

type Comments struct {
	ctx *service.Service
}

func (cm *Comments) Setup(ctx *service.Service) {
	cm.ctx = ctx

	// configure router
	CommentsRouter := cm.ctx.Echo.Group("/api/v2/comments")

	// auth middleware
	// CommentsRouter.Use()

	// CommentsRouter.GET("", cm.GetCommentsHandler)
	CommentsRouter.POST("", cm.CreateCommentHandler)
	// CommentsRouter.GET("/:id", cm.GetCommentHandler)
	// CommentsRouter.PUT("/:id", cm.UpdateCommentHandler)
	// CommentsRouter.DELETE("/:id", cm.DeleteCommentHandler)
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
