package comments

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Represent output data of GetCommentHandler
type GetCommentHandlerResponseBody struct {
	Comment *Comment `json:"comment" xml:"comment"`
	Message string   `json:"message" xml:"message"`
}

// Gets post by provided id from database, returns posts
func (c *Comments) GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	logger := c.ctx.Logger.Named("GetCommentHandler")

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

	// parse uuid
	logger.Infow("parsing uuid from path")
	uid, err := uuid.Parse(id)
	if err != nil {
		logger.Errorw("failed to parse uuid", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("uid", uid)

	// retreive comment from database
	logger.Infow("getting comment from database")
	var comment Comment
	err = c.ctx.MySQL.Model(&Comment{}).Where(&Comment{ID: uid}).First(&comment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Errorw("failed to find comment with provided id in database", "err", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		logger.Errorw("failed to get comment from database", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger = logger.With("comment", comment)

	// assemble response body
	logger.Infow("assembling response body")
	res := GetCommentHandlerResponseBody{
		Comment: &comment,
		Message: "successfully retrieved post",
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

	logger.Infow("successfully retrieved post by id from database")
	w.Write(b)
}
