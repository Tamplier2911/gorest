package comments

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

// Represent output data of GetCommentHandler
type GetCommentHandlerResponseBody struct {
	Comment *models.Comment `json:"comment" xml:"comment"`
	Message string          `json:"message" xml:"message"`
}

// Gets comment by provided id from database, returns comment
func (cm *Comments) GetCommentHandler(c echo.Context) error {
	logger := cm.ctx.Logger.Named("GetCommentHandler")

	// get id from path param
	logger.Infow("getting id from path params")
	id := c.Param("id")
	logger = logger.With("id", id)

	// parse uuid
	logger.Infow("parsing uuid from path")
	commentId, err := uuid.Parse(id)
	if err != nil {
		logger.Errorw("failed to parse uuid", "err", err)
		return cm.ResponseWriter(c, http.StatusBadRequest, GetCommentHandlerResponseBody{
			Message: "failed to parse uuid",
		})
	}
	logger = logger.With("commentId", commentId)

	// retreive comment from database
	logger.Infow("getting comment from database")
	var comment models.Comment
	err = cm.ctx.MySQL.
		Model(&models.Comment{}).
		Where(&models.Comment{Base: models.Base{ID: commentId}}).
		First(&comment).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Errorw("failed to find comment with provided id in database", "err", err)
			return cm.ResponseWriter(c, http.StatusBadRequest, GetCommentHandlerResponseBody{
				Message: "failed to find comment with provided id in database",
			})
		}

		logger.Errorw("failed to get comment from database", "err", err)
		return cm.ResponseWriter(c, http.StatusInternalServerError, GetCommentHandlerResponseBody{
			Message: "failed to get comment from database",
		})
	}
	logger = logger.With("comment", comment)

	// assemble response body
	logger.Infow("assembling response body")
	res := GetCommentHandlerResponseBody{
		Comment: &comment,
		Message: "successfully retrieved comment",
	}
	logger = logger.With("res", res)

	logger.Infow("successfully retrieved comment by id from database")
	return cm.ResponseWriter(c, http.StatusOK, res)
}
