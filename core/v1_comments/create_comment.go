package v1_comments

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
)

// Represent input data of CreateCommentHandler
type CreateCommentRequestBody struct {
	PostID string `json:"postId" form:"postId" url:"postId" binding:"required"`
	UserID string `json:"userId" form:"userId" url:"userId" binding:"required"`
	Name   string `json:"name" form:"name" url:"name" binding:"required"`
	Body   string `json:"body" form:"body" url:"body" binding:"required"`
}

// Represent output data of CreateCommentHandler
type CreateCommentResponseBody struct {
	Comment *models.Comment `json:"comment" xml:"comment"`
	Message string          `json:"message" xml:"message"`
}

// Creates comment instance and stores it in database
func (c *Comments) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	logger := c.Logger.Named("CreateCommentHandler")

	// parse body data
	logger.Infow("parsing request body")
	var body CreateCommentRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.Errorw("failed to parse request body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("body", body)

	// parse uuid id
	logger.Infow("parsing uuids from body")
	postUuid, err := uuid.Parse(string(body.PostID))
	if err != nil {
		logger.Errorw("failed to parse uuids from body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("postUuid", postUuid)

	userUuid, err := uuid.Parse(string(body.UserID))
	if err != nil {
		logger.Errorw("failed to parse uuids from body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("userUuid", userUuid)

	// save instance of comment in database
	logger.Infow("saving comment to database")
	comment := models.Comment{
		UserID: userUuid,
		PostID: postUuid,
		Name:   body.Name,
		Body:   body.Body,
	}
	err = c.MySQL.Model(&models.Comment{}).Create(&comment).Error
	if err != nil {
		logger.Errorw("failed to save comment in database", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := CreateCommentResponseBody{
		Comment: &comment,
		Message: "successfully created comment",
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
	w.WriteHeader(http.StatusCreated)

	logger.Debugw("successfully created comment record in database")
	w.Write(b)
}
