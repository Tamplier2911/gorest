package main

import (
	"encoding/json"
	"net/http"
)

// Represent output data of GetPostsHandler
type GetPostsHandlerResponseBody struct {
	Posts   []Post `json:"posts" xml:"posts"`
	Total   int64  `json:"total" xml:"total"`
	Message string `json:"message" xml:"message"`
}

// Get all posts from database, takes limit and offset query parameters, returns posts
func (s *Monolith) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	logger := s.Logger.Named("GetPostsHandler")

	// TODO: define limits and offsets with query params
	// retreive posts from database
	logger.Infow("getting posts from database")
	var total int64
	var posts []Post
	err := s.MySQL.Model(&Post{}).Count(&total).Limit(10).Offset(0).Find(&posts).Error
	if err != nil {
		logger.Errorw("failed to get posts from dataabase", "err", err)
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

	logger.Infow("successfully return all posts")
	w.Write(b)
}
