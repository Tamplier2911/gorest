package comments

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

	// get token from context
	token := access.GetTokenFromContext(c)
	logger = logger.With("token", token)

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

	// getting comment from database
	var comment models.Comment
	logger.Infow("getting comment from database")
	err = cm.MySQL.
		Model(&models.Comment{}).
		Where(&models.Comment{Base: models.Base{ID: commentId}}).
		First(&comment).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return cm.ResponseWriter(c, http.StatusNotFound, DeleteCommentHandlerResponseBody{
				Message: "failed to find record with provided id",
			})
		}
		logger.Errorw("failed to find comment record in database", "err", err)
		return cm.ResponseWriter(c, http.StatusInternalServerError, DeleteCommentHandlerResponseBody{
			Message: "failed to update comment",
		})
	}
	logger = logger.With("comment", comment)

	// check if user is comment author
	logger.Infow("checking if user is author a comment")
	if token.UserID != comment.UserID {
		logger.Errorw("user is not author of current comment", "err", err)
		return cm.ResponseWriter(c, http.StatusForbidden, DeleteCommentHandlerResponseBody{
			Message: "only author can change comment content",
		})
	}

	// delete comment from database
	logger.Infow("deleting comment from database")
	result := cm.MySQL.Delete(&comment)
	if result.Error != nil || result.RowsAffected == 0 {
		logger.Errorw("failed to delete comment with provided id from database", "err", err)
		return cm.ResponseWriter(c, http.StatusInternalServerError, DeleteCommentHandlerResponseBody{
			Message: "failed to delete comment",
		})
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := DeleteCommentHandlerResponseBody{
		Message: "successfully deleted comment",
	}
	logger = logger.With("res", res)

	logger.Infow("successfully deleted comment from database")
	return cm.ResponseWriter(c, http.StatusNoContent, res)
}
