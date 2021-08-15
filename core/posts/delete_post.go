package posts

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Represent output data of DeletePostHandler
type DeletePostHandlerResponseBody struct {
	Message string `json:"message" xml:"message"`
} // @name DeletePostResponse

// DeletePostHandler godoc
//
// @id				DeletePost
// @Summary 		Deletes post record.
// @Description 	Deletes post record from database using provided id.
//
// @Produce json
// @Produce xml
//
// @Success 200 	{object} DeletePostHandlerResponseBody
// @Failure 400,404 {object} DeletePostHandlerResponseBody
// @Failure 500 	{object} DeletePostHandlerResponseBody
// @Failure default {object} DeletePostHandlerResponseBody
//
// @Router /posts/{id} [DELETE]
func (p *Posts) DeletePostHandler(c echo.Context) error {
	logger := p.ctx.Logger.Named("DeletePostsHandler")

	// get id from path param
	logger.Infow("getting id from path params")
	id := c.Param("id")
	logger = logger.With("id", id)

	// parse uuid
	logger.Infow("parsing uuid from path")
	postId, err := uuid.Parse(id)
	if err != nil {
		logger.Errorw("failed to parse uuid", "err", err)
		return p.ResponseWriter(c, http.StatusBadRequest, DeletePostHandlerResponseBody{
			Message: "failed to parse uuid",
		})
	}
	logger = logger.With("postId", postId)

	// delete post from database
	logger.Infow("deleting post from database")
	result := p.ctx.MySQL.Model(&models.Post{}).Delete(&models.Post{Base: models.Base{ID: postId}})
	if result.Error != nil || result.RowsAffected == 0 {
		if result.Error == nil {
			return p.ResponseWriter(c, http.StatusBadRequest, DeletePostHandlerResponseBody{
				Message: "failed to find record with provided id",
			})
		}
		logger.Errorw("failed to delete post record from database", "err", err)
		return p.ResponseWriter(c, http.StatusInternalServerError, DeletePostHandlerResponseBody{
			Message: "failed to delete post from database",
		})
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := DeletePostHandlerResponseBody{
		Message: "successfully deleted post from database",
	}
	logger = logger.With("res", res)

	logger.Infow("successfully deleted post from database")
	return p.ResponseWriter(c, http.StatusNoContent, DeletePostHandlerResponseBody{
		Message: "successfully deleted post from database",
	})
}
