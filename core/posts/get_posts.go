package posts

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"

	"gorm.io/gorm/clause"
)

// Represent output data of GetPostsHandler
type GetPostsHandlerResponseBody struct {
	Posts   []Post `json:"posts" xml:"posts"`
	Total   int64  `json:"total" xml:"total"`
	Message string `json:"message" xml:"message"`
}

// Get all posts from database, takes limit and offset query parameters, returns posts
func (p *Posts) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	logger := p.ctx.Logger.Named("GetPostsHandler")

	// define db statement
	stmt := p.ctx.MySQL.Model(&Post{})

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

	// retreive posts from database
	logger.Infow("getting posts from database")
	var total int64
	var posts []Post
	err := stmt.
		Count(&total).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
		Find(&posts).
		Error
	if err != nil {
		logger.Errorw("failed to get posts from database", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger = logger.With("posts", posts)

	// assemble response body
	logger.Infow("assembling response body")
	res := GetPostsHandlerResponseBody{
		Posts:   posts,
		Total:   total,
		Message: "successfully retrieved posts",
	}
	logger = logger.With("res", res)

	// get clients accept header
	accept := r.Header.Get("Accept")

	var b []byte
	switch accept {
	case string(MimeTypesXML):
		// response with xml
		logger.Infow("marshaling response body to xml")
		w.Header().Set("Content-Type", string(MimeTypesXML))
		b, err = xml.Marshal(res)
	default:
		// default response with json
		logger.Infow("marshaling response body to json")
		w.Header().Set("Content-Type", string(MimeTypesJSON))
		b, err = json.Marshal(res)
	}

	if err != nil {
		logger.Errorw("failed to marshal response body", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write headers
	w.WriteHeader(http.StatusOK)

	logger.Infow("successfully retrieved all posts from database")
	w.Write(b)
}
