package main

import (
	"encoding/json"
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
func (s *Monolith) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.Named("UpdatePostHandler")

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
	result := s.MySQL.
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

	// response with xml or json based on condition

	// marshal json into bytes
	logger.Infow("marshaling response body")
	b, err := json.Marshal(res)
	if err != nil {
		logger.Errorw("failed to marshal response body", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write headers based
	w.Header().Set("Content-Type", "application/json")

	logger.Debugw("successfully updated post in database")
	w.Write(b)
}
