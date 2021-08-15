package posts

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// Represent input data of UpdatePostHandler
type UpdatePostRequestBody struct {
	Title string `json:"title" form:"title" binding:"required"`
	Body  string `json:"body" form:"body" binding:"required"`
}

// Represent output data of UpdatePostHandler
type UpdatePostResponseBody struct {
	Message string `json:"message" xml:"message"`
}

// Updates post instance in database
func (p *Posts) UpdatePostHandler(c echo.Context) error {
	logger := p.ctx.Logger.Named("UpdatePostHandler")

	// get id from path param
	logger.Infow("getting id from path params")
	id := c.Param("id")
	logger = logger.With("id", id)

	// parse uuid id
	logger.Infow("parsing uuid from path")
	uid, err := uuid.Parse(id)
	if err != nil {
		logger.Errorw("failed to parse uuid", "err", err)
		return p.ResponseWriter(c, http.StatusBadRequest, UpdatePostResponseBody{
			Message: "failed to parse uuid",
		})
	}
	logger = logger.With("uid", uid)

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

	// update post in database
	logger.Infow("updating post in database")
	result := p.ctx.MySQL.
		Model(&models.Post{}).
		Where(&models.Post{Base: models.Base{ID: uid}}).
		Updates(&models.Post{Title: body.Title, Body: body.Body})
	if result.Error != nil || result.RowsAffected == 0 {
		if result.Error == nil {
			return p.ResponseWriter(c, http.StatusInternalServerError, UpdatePostResponseBody{
				Message: "could not find record with this id to update",
			})
		}
		logger.Errorw("failed to update post in database", "err", result.Error)
		return p.ResponseWriter(c, http.StatusInternalServerError, UpdatePostResponseBody{
			Message: "failed to update post in database",
		})
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := UpdatePostResponseBody{
		Message: "successfully updated post",
	}
	logger = logger.With("res", res)

	logger.Debugw("successfully updated post in database")
	return p.ResponseWriter(c, http.StatusOK, UpdatePostResponseBody{
		Message: "successfully updated post in database",
	})
}
