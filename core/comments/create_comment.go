package comments

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Represent input data of CreateCommentHandler
type CreateCommentRequestBody struct {
	PostID string `json:"postId" form:"postId" binding:"required"`
	Name   string `json:"name" form:"name" binding:"required"`
	Body   string `json:"body" form:"body" binding:"required"`
} // @name CreateCommentRequest

// Represent output data of CreateCommentHandler
type CreateCommentResponseBody struct {
	Comment *models.Comment `json:"comment" xml:"comment"`
	Message string          `json:"message" xml:"message"`
} // @name CreateCommentResponse

// CreateCommentHandler godoc
//
// @id				CreateComment
// @Summary 		Creates comment record.
// @Description 	Creates comment record in database using provided data.
//
// @Tags			Comments
//
// @Accept json
//
// @Produce json
// @Produce xml
//
// @Param fields body CreateCommentRequestBody true "data"
//
// @Success 201 	{object} CreateCommentResponseBody
// @Failure 400,404 {object} CreateCommentResponseBody
// @Failure 500 	{object} CreateCommentResponseBody
// @Failure default {object} CreateCommentResponseBody
//
// @Router /comments [POST]
func (cm *Comments) CreateCommentHandler(c echo.Context) error {
	logger := cm.Logger.Named("CreateCommentHandler")

	// get token from context
	token := access.GetTokenFromContext(c)
	logger = logger.With("token", token)

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

	// save instance of comment in database
	logger.Infow("saving comment to database")
	comment := models.Comment{
		UserID: token.UserID,
		PostID: postUuid,
		Name:   body.Name,
		Body:   body.Body,
	}
	err = cm.MySQL.
		Model(&models.Comment{}).
		Create(&comment).
		Error
	if err != nil {
		logger.Errorw("failed to save comment in database", "err", err)
		return cm.ResponseWriter(c, http.StatusInternalServerError, CreateCommentResponseBody{
			Message: "failed to save comment",
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
