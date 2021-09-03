package comments

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// Represent output data of GetCommentsHandler
type GetCommentsHandlerResponseBody struct {
	Comments *[]models.Comment `json:"comments" xml:"comments"`
	Total    int64             `json:"total" xml:"total"`
	Message  string            `json:"message" xml:"message"`
} // @name GetCommentsResponse

// Get all comments from database, takes limit and offset query parameters, returns comments
func (c *Comments) GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	logger := c.Logger.Named("GetCommentsHandler")

	// define db statement
	stmt := c.MySQL.Model(&models.Comment{})

	limit := 10
	// get limit from query parameters
	lmt := r.FormValue("limit")
	if lmt != "" {
		logger.Infow("parsing limit query")
		lmtInt, err := strconv.Atoi(lmt)
		if err != nil {
			logger.Infow("invalid limit in query")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger = logger.With("limit", limit)
		limit = lmtInt
	}

	offset := 10
	// get offset from query parameters
	ofst := r.FormValue("offset")
	if ofst != "" {
		logger.Infow("parsing offset query")
		ofstInt, err := strconv.Atoi(ofst)
		if err != nil {
			logger.Infow("invalid offset in query")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger = logger.With("offset", offset)
		offset = ofstInt
	}

	// get post id from query parameters
	postId := r.FormValue("postId")
	if postId != "" {
		logger.Infow("parsing post id query")
		postUuid, err := uuid.Parse(postId)
		if err != nil {
			logger.Errorw("failed to parse uuid from body", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// add clause to statement
		stmt.Where(&models.Comment{PostID: postUuid})
	}

	// get user id from query parameters
	userId := r.FormValue("userId")
	if userId != "" {
		logger.Infow("parsing user id query")
		userUuid, err := uuid.Parse(userId)
		if err != nil {
			logger.Errorw("failed to parse uuid from body", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// add clause to statement
		stmt.Where(&models.Comment{UserID: userUuid})
	}

	// retreive comments from database
	logger.Infow("getting comments from database")
	var total int64
	var comments []models.Comment
	err := stmt.
		Count(&total).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
		Limit(limit).
		Offset(offset).
		Find(&comments).
		Error
	if err != nil {
		logger.Errorw("failed to get comments from database", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	// get clients accept header
	accept := r.Header.Get("Accept")

	var b []byte
	switch accept {
	case string(models.MimeTypesXML):
		// response with xml
		logger.Infow("marshaling response body to xml")
		w.Header().Set("Content-Type", string(models.MimeTypesXML))
		b, err = xml.Marshal(res)
	default:
		// default response with json
		logger.Infow("marshaling response body to json")
		w.Header().Set("Content-Type", string(models.MimeTypesJSON))
		b, err = json.Marshal(res)
	}

	if err != nil {
		logger.Errorw("failed to marshal response body", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write headers
	w.WriteHeader(http.StatusOK)

	logger.Infow("successfully retrieved all comments from database")
	w.Write(b)
}
