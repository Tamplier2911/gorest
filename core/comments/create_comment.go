package comments

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// Represent input data of CreateCommentHandler
type CreateCommentRequestBody struct {
	PostID string `json:"postId" form:"postId" binding:"required"`
	// TODO: consider getting user id from auth middleware in future
	UserID string `json:"userId" form:"userId" binding:"required"`
	Name   string `json:"name" form:"name" binding:"required"`
	Body   string `json:"body" form:"body" binding:"required"`
}

// Represent output data of CreateCommentHandler
type CreateCommentResponseBody struct {
	Comment *models.Comment `json:"comment" xml:"comment"`
	Message string          `json:"message" xml:"message"`
}

// Creates comment instance and stores it in database
func (cm *Comments) CreateCommentHandler(c echo.Context) error {
	logger := cm.ctx.Logger.Named("CreateCommentHandler")

	// parse body data
	logger.Infow("parsing request body")
	var body CreateCommentRequestBody
	err := c.Bind(&body)
	if err != nil {
		logger.Errorw("failed to parse request body", "err", err)
		return cm.ResponseWriter(c, http.StatusBadRequest, CreateCommentResponseBody{
			Message: "failed to parse request body",
		})
	}
	logger = logger.With("body", body)

	// parse uuids id
	logger.Infow("parsing uuids from body")
	postUuid, err := uuid.Parse(body.PostID)
	if err != nil {
		logger.Errorw("failed to parse uuids from body", "err", err)
		return cm.ResponseWriter(c, http.StatusBadRequest, CreateCommentResponseBody{
			Message: "failed to parse uuids from body",
		})
	}
	logger = logger.With("postUuid", postUuid)

	// TODO: consider getting user id from auth middleware in future
	userUuid, err := uuid.Parse(body.UserID)
	if err != nil {
		logger.Errorw("failed to parse uuids from body", "err", err)
		return cm.ResponseWriter(c, http.StatusBadRequest, CreateCommentResponseBody{
			Message: "failed to parse uuids from body",
		})
	}
	logger = logger.With("userUuid", userUuid)

	// save instance of comment in database
	logger.Infow("saving comment to database")
	comment := models.Comment{
		UserID: userUuid,
		PostID: postUuid,
		Name:   body.Name,
		Body:   body.Body,
	}
	err = cm.ctx.MySQL.Model(&models.Comment{}).Create(&comment).Error
	if err != nil {
		logger.Errorw("failed to save comment in database", "err", err)
		return cm.ResponseWriter(c, http.StatusInternalServerError, CreateCommentResponseBody{
			Message: "failed to save comment in database",
		})
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := CreateCommentResponseBody{
		Comment: &comment,
		Message: "successfully created comment",
	}
	logger = logger.With("res", res)

	logger.Debugw("successfully created comment record in database")
	return cm.ResponseWriter(c, http.StatusCreated, res)
}
