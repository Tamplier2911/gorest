package comments

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/Tamplier2911/gorest/pkg/models"
	"gorm.io/gorm/clause"
)

// Represent output data of GetCommentsHandler
type GetCommentsHandlerResponseBody struct {
	Comments []models.Comment `json:"comments" xml:"comments"`
	Total    int64            `json:"total" xml:"total"`
	Message  string           `json:"message" xml:"message"`
}

// Get all comments from database, takes limit and offset query parameters, returns comments
func (c *Comments) GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	logger := c.ctx.Logger.Named("GetCommentsHandler")

	// define db statement
	stmt := c.ctx.MySQL.Model(&models.Comment{})

	// TODO: consider refactoring that

	// get limit from query parameters
	limit := r.FormValue("limit")
	if limit != "" {
		logger.Infow("parsing limit query")
		lm, err := strconv.Atoi(limit)
		if err != nil {
			logger.Infow("invalid limit query")
			stmt.Limit(10)
		} else {
			stmt.Limit(lm)
			logger = logger.With("limit", lm)
		}
	}

	// get offset from query parameters
	offset := r.FormValue("offset")
	if offset != "" {
		logger.Infow("parsing offset query")
		of, err := strconv.Atoi(offset)
		if err != nil {
			logger.Infow("invalid offset query")
			stmt.Offset(0)
		} else {
			stmt.Offset(of)
			logger = logger.With("offset", of)
		}
	}

	// retreive comments from database
	logger.Infow("getting comments from database")
	var total int64
	var comments []models.Comment
	err := stmt.
		Count(&total).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
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
		Comments: comments,
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
