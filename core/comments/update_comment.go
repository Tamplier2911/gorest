package comments

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Represent input data of UpdateCommentHandler
type UpdateCommentRequestBody struct {
	Name string `json:"name" form:"name" binding:"required"`
	Body string `json:"body" form:"body" binding:"required"`
} // @name UpdateCommentRequest

// Represent output data of UpdateCommentHandler
type UpdateCommentResponseBody struct {
	Comment *models.Comment `json:"comment" xml:"comment"`
	Message string          `json:"message" xml:"message"`
} // @name UpdateCommentResponse

// UpdateCommentHandler godoc
//
// @id				UpdateComment
// @Summary 		Updates comment record.
// @Description 	Updates comment record in database using provided data.
//
// @Tags			Comments
//
// @Accept json
//
// @Produce json
// @Produce xml
//
// @Param fields body UpdateCommentRequestBody true "data"
//
// @Success 200 	{object} UpdateCommentResponseBody
// @Failure 400,404 {object} UpdateCommentResponseBody
// @Failure 500 	{object} UpdateCommentResponseBody
// @Failure default {object} UpdateCommentResponseBody
//
// @Router /comments/{id} [PUT]
func (cm *Comments) UpdateCommentHandler(c echo.Context) error {
	logger := cm.Logger.Named("UpdateCommentHandler")

	// get token from context
	token := access.GetTokenFromContext(c)
	logger = logger.With("token", token)

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
			return cm.ResponseWriter(c, http.StatusNotFound, UpdateCommentResponseBody{
				Message: "failed to find record with provided id",
			})
		}
		logger.Errorw("failed to find comment record in database", "err", err)
		return cm.ResponseWriter(c, http.StatusInternalServerError, UpdateCommentResponseBody{
			Message: "failed to update comment",
		})
	}
	logger = logger.With("comment", comment)

	// check if user is comment author
	logger.Infow("checking if user is author a comment")
	if token.UserID != comment.UserID {
		logger.Errorw("user is not author of current comment", "err", err)
		return cm.ResponseWriter(c, http.StatusForbidden, UpdateCommentResponseBody{
			Message: "only author can change comment content",
		})
	}

	// update comment in database
	logger.Infow("updating comment in database")
	err = cm.MySQL.
		Model(&comment).
		Updates(&models.Comment{Name: body.Name, Body: body.Body}).
		Error
	if err != nil {
		logger.Errorw("failed to update comment in database", "err", err)
		return cm.ResponseWriter(c, http.StatusInternalServerError, UpdateCommentResponseBody{
			Message: "failed to update comment",
		})
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := UpdateCommentResponseBody{
		Comment: &comment,
		Message: "successfully updated comment",
	}
	logger = logger.With("res", res)

	logger.Debugw("successfully updated comment in database")
	return cm.ResponseWriter(c, http.StatusOK, res)
}
