package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// Represent input data of CreatePostHandler
type CreatePostRequestBody struct {
	UserID string `json:"userId" form:"userId" url:"userId" binding:"required"`
	Title  string `json:"title" form:"title" url:"title" binding:"required"`
	Body   string `json:"body" form:"body" url:"body" binding:"required"`
}

// Represent output data of CreatePostHandler
type CreatePostResponseBody struct {
	Post    *Post  `json:"post"`
	Message string `json:"message"`
}

// Creates post instance and stores it in database
func (s *Monolith) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.Named("CreatePostHandler")

	// parse body data
	logger.Infow("parsing request body")
	var req CreatePostRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Errorw("failed to parse request body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("req", req)

	// parse uuid id
	logger.Infow("parsing uuid from path")
	uid, err := uuid.Parse(string(req.UserID))
	if err != nil {
		logger.Errorw("failed to parse", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("uid", uid)

	// save instance of post in database
	logger.Infow("saving post to database")
	post := Post{
		UserID: uid,
		Title:  req.Title,
		Body:   req.Body,
	}
	err = s.MySQL.Model(&Post{}).Create(&post).Error
	if err != nil {
		logger.Errorw("failed to save post in database", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := CreatePostResponseBody{
		Post:    &post,
		Message: "successfully created post",
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

	logger.Debugw("successfully return all posts")
	w.Write(b)
}
