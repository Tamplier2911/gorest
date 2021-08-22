package posts

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
)

// Represent input query of GetPostHandler
type GetPostsHandlerRequestQuery struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
} // @name GetPostRequest

// Represent output data of GetPostsHandler
type GetPostsHandlerResponseBody struct {
	Posts   *[]models.Post `json:"posts" xml:"posts"`
	Total   int64          `json:"total" xml:"total"`
	Message string         `json:"message" xml:"message"`
} // @name GetPostResponse

// GetPostsHandler godoc
//
// @id				GetPosts
// @Summary 		Gets post records.
// @Description 	Gets post records from database using provided query.
//
// @Tags			Posts
//
// @Produce json
// @Produce xml
//
// @Param fields query GetPostsHandlerRequestQuery true "data"
//
// @Success 200 	{object} GetPostsHandlerResponseBody
// @Failure 400,404 {object} GetPostsHandlerResponseBody
// @Failure 500 	{object} GetPostsHandlerResponseBody
// @Failure default {object} GetPostsHandlerResponseBody
//
// @Router /posts [GET]
func (p *Posts) GetPostsHandler(c echo.Context) error {
	logger := p.Logger.Named("GetPostsHandler")

	logger.Infow("parsing request query params")
	var query GetPostsHandlerRequestQuery
	err := c.Bind(&query)
	if err != nil {
		logger.Errorw("failed to parse request body", "err", err)
		return p.ResponseWriter(c, http.StatusBadRequest, CreatePostResponseBody{
			Message: "failed to parse request body",
		})
	}
	logger = logger.With("query", query)

	// set default limit to 10
	limit := 10
	if query.Limit != 0 {
		limit = query.Limit
	}

	// retreive posts from database
	logger.Infow("getting posts from database")
	var total int64
	var posts []models.Post
	err = p.MySQL.Model(&models.Post{}).
		Count(&total).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
		Limit(limit).
		Offset(query.Offset).
		Find(&posts).
		Error
	if err != nil {
		logger.Errorw("failed to get posts from database", "err", err)
		return p.ResponseWriter(c, http.StatusInternalServerError, GetPostsHandlerResponseBody{
			Message: "failed to get posts",
		})
	}
	logger = logger.With("posts", posts)

	// assemble response body
	logger.Infow("assembling response body")
	res := GetPostsHandlerResponseBody{
		Posts:   &posts,
		Total:   total,
		Message: "successfully retrieved posts",
	}
	logger = logger.With("res", res)

	logger.Infow("successfully retrieved all posts from database")
	return p.ResponseWriter(c, http.StatusOK, res)
}
