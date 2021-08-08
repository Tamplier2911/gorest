package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Represent output data of GetPostHandler
type GetPostHandlerResponseBody struct {
	Post    *Post  `json:"post" xml:"posts"`
	Message string `json:"message" xml:"message"`
}

// Gets post by provided id from database, returns posts
func (s *Monolith) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.Named("GetPostHandler")

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

	// retreive post from database
	logger.Infow("getting post from database")
	var post Post
	err = s.MySQL.Model(&Post{}).Where(&Post{ID: uid}).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Errorw("failed to find post with provided id in database", "err", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		logger.Errorw("failed to get posts from database", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger = logger.With("post", post)

	// assemble response body
	logger.Infow("assembling response body")
	res := GetPostHandlerResponseBody{
		Post:    &post,
		Message: "successfully retrieved post",
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
	w.WriteHeader(http.StatusOK)

	logger.Infow("successfully retrieved post by id from database")
	w.Write(b)
}
