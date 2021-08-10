package posts

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// Represent input data of UpdatePostHandler
type UpdatePostRequestBody struct {
	Title string `json:"title" form:"title" url:"title" binding:"required"`
	Body  string `json:"body" form:"body" url:"body" binding:"required"`
}

// Represent output data of UpdatePostHandler
type UpdatePostResponseBody struct {
	Message string `json:"message"`
}

// Updates post instance in database
func (p *Posts) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	logger := p.ctx.Logger.Named("UpdatePostHandler")

	// TODO: consider abstracting this to a middleware
	// get id from path
	logger.Infow("getting id from path")
	pathSlice := strings.Split(r.URL.Path, "/")
	id := pathSlice[len(pathSlice)-1]
	if id == "" {
		err := errors.New("failed to get id from path")
		logger.Errorw("failed to get id from path", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("id", id)

	// parse uuid id
	logger.Infow("parsing uuid from path")
	uid, err := uuid.Parse(id)
	if err != nil {
		logger.Errorw("failed to parse", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("uid", uid)

	// parse body data
	logger.Infow("parsing request body")
	var body CreatePostRequestBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.Errorw("failed to parse request body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("req", body)

	// update post in database
	logger.Infow("updating post in database")
	result := p.ctx.MySQL.
		Model(&Post{}).
		Where(&Post{ID: uid}).
		Updates(&Post{Title: body.Title, Body: body.Body})
	if result.Error != nil || result.RowsAffected == 0 {
		if result.Error == nil {
			result.Error = errors.New("record not found")
		}
		logger.Errorw("failed to update post in database", "err", err)
		http.Error(w, result.Error.Error(), http.StatusBadRequest)
		return
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := UpdatePostResponseBody{
		Message: "successfully updated post",
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

	logger.Debugw("successfully updated post in database")
	w.Write(b)
}