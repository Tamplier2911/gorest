package comments

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
)

// Represent intput data of GetCommentsHandler
type GetCommentsHandlerRequestQuery struct {
	Limit  int    `query:"limit"`
	Offset int    `query:"offset"`
	UserID string `query:"userId"`
	PostID string `query:"postId"`
} // @name GetCommentsRequest

// Represent output data of GetCommentsHandler
type GetCommentsHandlerResponseBody struct {
	Comments *[]models.Comment `json:"comments" xml:"comments"`
	Total    int64             `json:"total" xml:"total"`
	Message  string            `json:"message" xml:"message"`
} // @name GetCommentsResponse

// GetCommentsHandler godoc
//
// @id				GetComments
// @Summary 		Gets comment records.
// @Description 	Gets comment records from database using provided query.
//
// @Tags			Comments
//
// @Produce json
// @Produce xml
//
// @Param fields query GetCommentsHandlerRequestQuery true "data"
//
// @Success 200 	{object} GetCommentsHandlerResponseBody
// @Failure 400,404 {object} GetCommentsHandlerResponseBody
// @Failure 500 	{object} GetCommentsHandlerResponseBody
// @Failure default {object} GetCommentsHandlerResponseBody
//
// @Router /comments [GET]
func (cm *Comments) GetCommentsHandler(c echo.Context) error {
	logger := cm.Logger.Named("GetCommentsHandler")

	logger.Infow("parsing request query params")
	var query GetCommentsHandlerRequestQuery
	err := c.Bind(&query)
	if err != nil {
		logger.Errorw("failed to parse request body", "err", err)
		return cm.ResponseWriter(c, http.StatusBadRequest, GetCommentsHandlerResponseBody{
			Message: "failed to parse request body",
		})
	}
	logger = logger.With("query", query)

	stmt := cm.MySQL.Model(&models.Comment{})

	// append post id to where clause
	if query.PostID != "" {
		logger.Infow("parsing uuid form query")
		postUuid, err := uuid.Parse(query.PostID)
		if err != nil {
			logger.Errorw("failed to parse uuids from body", "err", err)
			return cm.ResponseWriter(c, http.StatusBadRequest, CreateCommentResponseBody{
				Message: "failed to parse uuids from body",
			})
		}
		logger = logger.With("postUuid", postUuid)

		// add clause to statement
		stmt.Where(&models.Comment{PostID: postUuid})
	}

	// append user id to where clause
	if query.UserID != "" {
		logger.Infow("parsing uuid form query")
		userUuid, err := uuid.Parse(query.UserID)
		if err != nil {
			logger.Errorw("failed to parse uuids from body", "err", err)
			return cm.ResponseWriter(c, http.StatusBadRequest, CreateCommentResponseBody{
				Message: "failed to parse uuids from body",
			})
		}
		logger = logger.With("userUuid", userUuid)

		// add clause to statement
		stmt.Where(&models.Comment{UserID: userUuid})
	}

	// set default limit to 10
	limit := 10
	if query.Limit != 0 {
		limit = query.Limit
	}

	// retreive comments from database
	logger.Infow("getting comments from database")
	var total int64
	var comments []models.Comment
	err = stmt.Count(&total).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
		Limit(limit).
		Offset(query.Offset).
		Find(&comments).
		Error
	if err != nil {
		logger.Errorw("failed to get comments from database", "err", err)
		return cm.ResponseWriter(c, http.StatusInternalServerError, GetCommentsHandlerResponseBody{
			Message: "failed to get comments",
		})
	}
	logger = logger.With("comments", comments)

	// assemble response body
	logger.Infow("assembling response body")
	res := GetCommentsHandlerResponseBody{
		Comments: &comments,
		Total:    total,
		Message:  "successfully retrieved comments",
	}
	logger = logger.With("res", res)

	logger.Infow("successfully retrieved all comments from database")
	return cm.ResponseWriter(c, http.StatusOK, res)
}
