package posts

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Represent input data of UpdatePostHandler
type UpdatePostRequestBody struct {
	Title string `json:"title" form:"title" binding:"required"`
	Body  string `json:"body" form:"body" binding:"required"`
} // @name UpdatePostRequest

// Represent output data of UpdatePostHandler
type UpdatePostResponseBody struct {
	Post    *models.Post `json:"post" xml:"post"`
	Message string       `json:"message" xml:"message"`
} // @name UpdatePostResponse

// UpdatePostHandler godoc
//
// @id				UpdatePost
// @Summary 		Updates post record.
// @Description 	Updates post record in database using provided data.
//
// @Tags			Posts
//
// @Accept json
//
// @Produce json
// @Produce xml
//
// @Param fields body UpdatePostRequestBody true "data"
//
// @Success 200 	{object} UpdatePostResponseBody
// @Failure 400,404 {object} UpdatePostResponseBody
// @Failure 500 	{object} UpdatePostResponseBody
// @Failure default {object} UpdatePostResponseBody
//
// @Router /posts/{id} [PUT]
func (p *Posts) UpdatePostHandler(c echo.Context) error {
	logger := p.Logger.Named("UpdatePostHandler")

	// get token from context
	token := access.GetTokenFromContext(c)
	logger = logger.With("token", token)

	// get id from path param
	logger.Infow("getting id from path params")
	id := c.Param("id")
	logger = logger.With("id", id)

	// parse uuid id
	logger.Infow("parsing uuid from path")
	postId, err := uuid.Parse(id)
	if err != nil {
		logger.Errorw("failed to parse uuid", "err", err)
		return p.ResponseWriter(c, http.StatusBadRequest, UpdatePostResponseBody{
			Message: "failed to parse uuid",
		})
	}
	logger = logger.With("postId", postId)

	// parse body data
	logger.Infow("parsing request body")
	var body CreatePostRequestBody
	err = c.Bind(&body)
	if err != nil {
		logger.Errorw("failed to parse request body", "err", err)
		return p.ResponseWriter(c, http.StatusBadRequest, UpdatePostResponseBody{
			Message: "failed to parse request body",
		})
	}
	logger = logger.With("body", body)

	// get post from database
	var post models.Post
	logger.Infow("getting post from database")
	err = p.MySQL.
		Model(&models.Post{}).
		Where(&models.Post{Base: models.Base{ID: postId}}).
		First(&post).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return p.ResponseWriter(c, http.StatusNotFound, UpdatePostResponseBody{
				Message: "failed to find record with provided id",
			})
		}
		logger.Errorw("failed to find post record in database", "err", err)
		return p.ResponseWriter(c, http.StatusInternalServerError, UpdatePostResponseBody{
			Message: "failed to update post",
		})
	}
	logger = logger.With("post", post)

	// check if user is post author
	logger.Infow("checking if user is author a post")
	if token.UserID != post.UserID {
		logger.Errorw("user is not author of current post", "err", err)
		return p.ResponseWriter(c, http.StatusForbidden, UpdatePostResponseBody{
			Message: "only author can change post content",
		})
	}

	// update post in database
	logger.Infow("updating post in database")
	err = p.MySQL.
		Model(&post).
		Updates(&models.Post{Title: body.Title, Body: body.Body}).
		Error
	if err != nil {
		logger.Errorw("failed to update post in database", "err", err)
		return p.ResponseWriter(c, http.StatusInternalServerError, UpdatePostResponseBody{
			Message: "failed to update post",
		})
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := UpdatePostResponseBody{
		Post:    &post,
		Message: "successfully updated post",
	}
	logger = logger.With("res", res)

	logger.Debugw("successfully updated post in database")
	return p.ResponseWriter(c, http.StatusOK, res)
}
