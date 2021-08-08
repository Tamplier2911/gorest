package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Represent output data of DeletePostHandler
type DeletePostHandlerResponseBody struct {
	Message string `json:"message" xml:"message"`
}

// Deletes post by provided id from database
func (s *Monolith) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.Named("DeletePostsHandler")

	// TODO: consider abstracting this to a middleware
	// get id from path
	logger.Infow("getting id from path")
	p := strings.Split(r.URL.Path, "/")
	id := p[len(p)-1]
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
		logger.Errorw("failed to parse", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("uid", uid)

	// delete post from database
	logger.Infow("deleting post from database")
	err = s.MySQL.Model(&Post{}).Delete(&Post{ID: uid}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Errorw("failed find post with provided id in database", "err", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		logger.Errorw("failed to delete post with provided id from database", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := DeletePostHandlerResponseBody{
		Message: "successfully deleted post from database",
	}
	logger = logger.With("res", res)

	// response with xml or json based on condition
	logger.Infow("marshaling response body")
	b, err := json.Marshal(res)
	if err != nil {
		logger.Errorw("failed to marshal response body", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write headers
	w.Header().Set("Content-Type", "application/json")

	logger.Infow("successfully deleted post from database")
	w.Write(b)
}
