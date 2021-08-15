package comments

import (
	"errors"
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// Represent input data of UpdateCommentHandler
type UpdateCommentRequestBody struct {
	Name string `json:"name" form:"name" binding:"required"`
	Body string `json:"body" form:"body" binding:"required"`
}

// Represent output data of UpdateCommentHandler
type UpdateCommentResponseBody struct {
	Message string `json:"message" xml:"message"`
}

// Updates post instance in database
func (cm *Comments) UpdateCommentHandler(c echo.Context) error {
	logger := cm.ctx.Logger.Named("UpdateCommentHandler")

	// get id from path param
	logger.Infow("getting id from path params")
	id := c.Param("id")
	logger = logger.With("id", id)

	// parse uuid id
	logger.Infow("parsing uuid from path")
	commentId, err := uuid.Parse(id)
	if err != nil {
		logger.Errorw("failed to parse uuid", "err", err)
		return cm.ResponseWriter(c, http.StatusBadRequest, UpdateCommentResponseBody{
			Message: "failed to parse uuid",
		})
	}
	logger = logger.With("commentId", commentId)

	// parse body data
	logger.Infow("parsing request body")
	var body CreateCommentRequestBody
	err = c.Bind(&body)
	if err != nil {
		logger.Errorw("failed to parse request body", "err", err)
		return cm.ResponseWriter(c, http.StatusBadRequest, UpdateCommentResponseBody{
			Message: "failed to parse request body",
		})
	}
	logger = logger.With("body", body)

	// update post in database
	logger.Infow("updating post in database")
	result := cm.ctx.MySQL.
		Model(&models.Comment{}).
		Where(&models.Comment{Base: models.Base{ID: commentId}}).
		Updates(&models.Comment{Name: body.Name, Body: body.Body})
	if result.Error != nil || result.RowsAffected == 0 {
		if result.Error == nil {
			result.Error = errors.New("record not found")
			return cm.ResponseWriter(c, http.StatusBadRequest, UpdateCommentResponseBody{
				Message: "failed to find comment with provided id in database",
			})
		}
		logger.Errorw("failed to update post in database", "err", err)
		return cm.ResponseWriter(c, http.StatusInternalServerError, UpdateCommentResponseBody{
			Message: "failed to update post in database",
		})
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := UpdateCommentResponseBody{
		Message: "successfully updated post",
	}
	logger = logger.With("res", res)

	logger.Debugw("successfully updated post in database")
	return cm.ResponseWriter(c, http.StatusOK, res)
}
