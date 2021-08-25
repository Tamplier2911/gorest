package posts

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
// @Tags			Posts
//
// @Produce json
// @Produce xml
//
// @Success 204 	{object} DeletePostHandlerResponseBody
// @Failure 400,404 {object} DeletePostHandlerResponseBody
// @Failure 500 	{object} DeletePostHandlerResponseBody
// @Failure default {object} DeletePostHandlerResponseBody
//
// @Security ApiKeyAuth
//
// @Router /posts/{id} [DELETE]
func (p *Posts) DeletePostHandler(c echo.Context) error {
	logger := p.Logger.Named("DeletePostsHandler")

	// get token from context
	token := access.GetTokenFromContext(c)
	logger = logger.With("token", token)

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
			return p.ResponseWriter(c, http.StatusNotFound, DeletePostHandlerResponseBody{
				Message: "failed to find record with provided id",
			})
		}
		logger.Errorw("failed to find post record in database", "err", err)
		return p.ResponseWriter(c, http.StatusInternalServerError, DeletePostHandlerResponseBody{
			Message: "failed to delete post",
		})
	}
	logger = logger.With("post", post)

	// check if user is post author
	logger.Infow("checking if user is author a post")
	if token.UserID != post.UserID {
		logger.Errorw("user is not author of current post", "err", err)
		return p.ResponseWriter(c, http.StatusForbidden, DeletePostHandlerResponseBody{
			Message: "only author can delete post",
		})
	}

	// delete post from database
	logger.Infow("deleting post from database")
	err = p.MySQL.
		Delete(&post).
		Error
	if err != nil {
		logger.Errorw("failed to delete post record from database", "err", err)
		return p.ResponseWriter(c, http.StatusInternalServerError, DeletePostHandlerResponseBody{
			Message: "failed to delete post",
		})
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := DeletePostHandlerResponseBody{
		Message: "successfully deleted post",
	}
	logger = logger.With("res", res)

	logger.Infow("successfully deleted post from database")
	return p.ResponseWriter(c, http.StatusNoContent, res)
}
