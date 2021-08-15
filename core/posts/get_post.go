package posts

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Represent output data of GetPostHandler
type GetPostHandlerResponseBody struct {
	Post    *models.Post `json:"post" xml:"posts"`
	Message string       `json:"message" xml:"message"`
} // @name GetPostResponse

// GetPostHandler godoc
//
// @id				GetPost
// @Summary 		Gets post record.
// @Description 	Gets post record from database using provided id.
//
// @Produce json
// @Produce xml
//
// @Success 200 	{object} GetPostHandlerResponseBody
// @Failure 400,404 {object} GetPostHandlerResponseBody
// @Failure 500 	{object} GetPostHandlerResponseBody
// @Failure default {object} GetPostHandlerResponseBody
//
// @Router /posts/{id} [GET]
func (p *Posts) GetPostHandler(c echo.Context) error {
	logger := p.ctx.Logger.Named("GetPostHandler")

	// get id from path param
	logger.Infow("getting id from path params")
	id := c.Param("id")
	logger = logger.With("id", id)

	// parse uuid
	logger.Infow("parsing uuid from path")
	postId, err := uuid.Parse(id)
	if err != nil {
		logger.Errorw("failed to parse uuid", "err", err)
		return p.ResponseWriter(c, http.StatusBadRequest, GetPostHandlerResponseBody{
			Message: "failed to parse uuid",
		})
	}
	logger = logger.With("postId", postId)

	// retreive post from database
	logger.Infow("getting post from database")
	var post models.Post
	err = p.ctx.MySQL.Model(&models.Post{}).Where(&models.Post{Base: models.Base{ID: postId}}).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Errorw("failed to find post with provided id in database", "err", err)
			return p.ResponseWriter(c, http.StatusNotFound, GetPostHandlerResponseBody{
				Message: "failed to find post with provided id in database",
			})
		}

		logger.Errorw("failed to get posts from database", "err", err)
		return p.ResponseWriter(c, http.StatusInternalServerError, GetPostHandlerResponseBody{
			Message: "failed to get posts from database",
		})
	}
	logger = logger.With("post", post)

	// assemble response body
	logger.Infow("assembling response body")
	res := GetPostHandlerResponseBody{
		Post:    &post,
		Message: "successfully retrieved post",
	}
	logger = logger.With("res", res)

	logger.Infow("successfully retrieved post by id from database")
	return p.ResponseWriter(c, http.StatusOK, res)
}
