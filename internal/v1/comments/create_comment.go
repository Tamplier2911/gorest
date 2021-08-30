package comments

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
)

// Represent input data of CreateCommentHandler
type CreateCommentHandlerRequestBody struct {
	PostID string `json:"postId" form:"postId" url:"postId" binding:"required" validate:"required"`
	Name   string `json:"name" form:"name" url:"name" binding:"required" validate:"required"`
	Body   string `json:"body" form:"body" url:"body" binding:"required" validate:"required"`
} // @name CreateCommentRequest

// Represent output data of CreateCommentHandler
type CreateCommentHandlerResponseBody struct {
	Comment *models.Comment `json:"comment" xml:"comment"`
	Message string          `json:"message" xml:"message"`
} // @name CreateCommentResponse

// Creates comment instance and stores it in database
func (c *Comments) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	logger := c.Logger.Named("CreateCommentHandler")

	// get token from context
	token := r.Context().Value("token").(*access.Token)
	logger = logger.With("token", token)

	// parse body data
	logger.Infow("parsing request body")
	var body CreateCommentHandlerRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.Errorw("failed to parse request body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("body", body)

	// validate body data
	logger.Infow("validating request body")
	err = c.Validator.Struct(&body)
	if err != nil {
		logger.Errorw("failed to validate body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// parse uuid id
	logger.Infow("parsing uuids from body")
	postUuid, err := uuid.Parse(string(body.PostID))
	if err != nil {
		logger.Errorw("failed to parse uuids from body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("postUuid", postUuid)

	// save instance of comment in database
	logger.Infow("saving comment to database")
	comment := models.Comment{
		UserID: token.UserID,
		PostID: postUuid,
		Name:   body.Name,
		Body:   body.Body,
	}
	err = c.MySQL.
		Model(&models.Comment{}).
		Create(&comment).
		Error
	if err != nil {
		logger.Errorw("failed to save comment in database", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := CreateCommentHandlerResponseBody{
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

	logger.Infow("successfully created comment record in database")
	w.Write(b)
}
