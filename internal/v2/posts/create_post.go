package posts

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/labstack/echo/v4"
)

// Represent input data of CreatePostHandler
type CreatePostRequestBody struct {
	Title string `json:"title" form:"title" binding:"required"`
	Body  string `json:"body" form:"body" binding:"required"`
} // @name CreatePostRequest

// Represent output data of CreatePostHandler
type CreatePostResponseBody struct {
	Post    *models.Post `json:"post" xml:"post"`
	Message string       `json:"message" xml:"message"`
} // @name CreatePostResponse

// CreatePostHandler godoc
//
// @id				CreatePost
// @Summary 		Creates post record.
// @Description 	Creates post record in database using provided data.
//
// @Tags			Posts
//
// @Accept json
//
// @Produce json
// @Produce xml
//
// @Param fields body CreatePostRequestBody true "data"
//
// @Success 201 	{object} CreatePostResponseBody
// @Failure 400,404 {object} CreatePostResponseBody
// @Failure 500 	{object} CreatePostResponseBody
// @Failure default {object} CreatePostResponseBody
//
// @Security ApiKeyAuth
//
// @Router /posts [POST]
func (p *Posts) CreatePostHandler(c echo.Context) error {
	logger := p.Logger.Named("CreatePostHandler")

	// get token from context
	token := access.GetTokenFromContext(c)
	logger = logger.With("token", token)

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

	// save instance of post in database
	logger.Infow("saving post to database")
	post := models.Post{
		UserID: token.UserID,
		Title:  body.Title,
		Body:   body.Body,
	}
	err = p.MySQL.Model(&models.Post{}).Create(&post).Error
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
	return p.ResponseWriter(c, http.StatusCreated, res)
}