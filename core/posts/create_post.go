package posts

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// Represent input data of CreatePostHandler
type CreatePostRequestBody struct {
	UserID string `json:"userId" form:"userId" url:"userId" query:"userId" binding:"required"`
	Title  string `json:"title" form:"title" url:"title" query:"userId" binding:"required"`
	Body   string `json:"body" form:"body" url:"body" query:"userId" binding:"required"`
}

// Represent output data of CreatePostHandler
type CreatePostResponseBody struct {
	Post    *models.Post `json:"post" xml:"post"`
	Message string       `json:"message" xml:"message"`
}

// Creates post instance and stores it in database
func (p *Posts) CreatePostHandler(c echo.Context) error {
	logger := p.ctx.Logger.Named("CreatePostHandler")

	// parse body data
	logger.Infow("parsing request body")
	var body CreatePostRequestBody
	err := c.Bind(&body)
	if err != nil {
		logger.Errorw("failed to parse request body", "err", err)
		return p.ResponseWriter(c, http.StatusBadRequest, CreatePostResponseBody{
			Message: "failed to parse request body",
		})
	}
	logger = logger.With("body", body)

	// parse uuid id
	logger.Infow("parsing uuid from body")
	uid, err := uuid.Parse(string(body.UserID))
	if err != nil {
		logger.Errorw("failed to parse uuid from body", "err", err)
		return p.ResponseWriter(c, http.StatusBadRequest, CreatePostResponseBody{
			Message: "failed to parse request uuid",
		})
	}
	logger = logger.With("uid", uid)

	// save instance of post in database
	logger.Infow("saving post to database")
	post := models.Post{
		UserID: uid,
		Title:  body.Title,
		Body:   body.Body,
	}
	err = p.ctx.MySQL.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		logger.Errorw("failed to save post in database", "err", err)
		return p.ResponseWriter(c, http.StatusInternalServerError, CreatePostResponseBody{
			Message: "failed to create post",
		})
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := CreatePostResponseBody{
		Post:    &post,
		Message: "successfully created post",
	}
	logger = logger.With("res", res)

	logger.Infow("successfully created post")
	return p.ResponseWriter(c, http.StatusOK, res)
}
