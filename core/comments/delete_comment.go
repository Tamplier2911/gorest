package comments

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Represent output data of DeleteCommentHandler
type DeleteCommentHandlerResponseBody struct {
	Message string `json:"message" xml:"message"`
} // @name DeleteCommentResponse

// DeleteCommentHandler godoc
//
// @id				DeleteComment
// @Summary 		Deletes comment record.
// @Description 	Deletes comment record from database using provided id.
//
// @Tags			Comments
//
// @Produce json
// @Produce xml
//
// @Success 204 	{object} DeleteCommentHandlerResponseBody
// @Failure 400,404 {object} DeleteCommentHandlerResponseBody
// @Failure 500 	{object} DeleteCommentHandlerResponseBody
// @Failure default {object} DeleteCommentHandlerResponseBody
//
// @Router /comments/{id} [DELETE]
func (cm *Comments) DeleteCommentHandler(c echo.Context) error {
	logger := cm.Logger.Named("DeleteCommentHandler")

	// get id from path param
	logger.Infow("getting id from path params")
	id := c.Param("id")
	logger = logger.With("id", id)

	// parse uuid
	logger.Infow("parsing uuid from path")
	commentId, err := uuid.Parse(id)
	if err != nil {
		logger.Errorw("failed to parse uuid", "err", err)
		return cm.ResponseWriter(c, http.StatusBadRequest, DeleteCommentHandlerResponseBody{
			Message: "failed to parse uuid",
		})
	}
	logger = logger.With("commentId", commentId)

	// delete comment from database
	logger.Infow("deleting comment from database")
	result := cm.MySQL.Model(&models.Comment{}).Delete(&models.Comment{Base: models.Base{ID: commentId}})
	if result.Error != nil || result.RowsAffected == 0 {
		if result.Error == nil {
			return cm.ResponseWriter(c, http.StatusNotFound, DeleteCommentHandlerResponseBody{
				Message: "failed to find comment with provided id in database",
			})
		}
		logger.Errorw("failed to delete comment with provided id from database", "err", err)
		return cm.ResponseWriter(c, http.StatusInternalServerError, DeleteCommentHandlerResponseBody{
			Message: "failed to delete comment with provided id from database",
		})
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := DeleteCommentHandlerResponseBody{
		Message: "successfully deleted comment from database",
	}
	logger = logger.With("res", res)

	logger.Infow("successfully deleted comment from database")
	return cm.ResponseWriter(c, http.StatusNoContent, res)
}
