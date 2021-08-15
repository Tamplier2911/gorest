package posts

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

// Represent output data of GetPostHandler
type GetPostHandlerResponseBody struct {
	Post    *models.Post `json:"post" xml:"posts"`
	Message string       `json:"message" xml:"message"`
}

// Gets post by provided id from database, returns posts
func (p *Posts) GetPostHandler(c echo.Context) error {
	logger := p.ctx.Logger.Named("GetPostHandler")

	// get id from path param
	logger.Infow("getting id from path params")
	id := c.Param("id")
	logger = logger.With("id", id)

	// parse uuid
	logger.Infow("parsing uuid from path")
	uid, err := uuid.Parse(id)
	if err != nil {
		logger.Errorw("failed to parse uuid", "err", err)
		return p.ResponseWriter(c, http.StatusBadRequest, GetPostHandlerResponseBody{
			Message: "failed to parse uuid",
		})
	}
	logger = logger.With("uid", uid)

	// retreive post from database
	logger.Infow("getting post from database")
	var post models.Post
	err = p.ctx.MySQL.Model(&models.Post{}).Where(&models.Post{Base: models.Base{ID: uid}}).First(&post).Error
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
