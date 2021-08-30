package posts

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"strings"

	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Represent input data of UpdatePostHandler
type UpdatePostHandlerRequestBody struct {
	Title string `json:"title" form:"title" url:"title" binding:"required" validate:"required"`
	Body  string `json:"body" form:"body" url:"body" binding:"required" validate:"required"`
} // @name UpdatePostRequest

// Represent output data of UpdatePostHandler
type UpdatePostHandlerResponseBody struct {
	Post    *models.Post `json:"post" xml:"post"`
	Message string       `json:"message" xml:"message"`
} // @name UpdatePostResponse

// Updates post instance in database
func (p *Posts) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	logger := p.Logger.Named("UpdatePostHandler")

	// get token from context
	token := r.Context().Value("token").(*access.Token)
	logger = logger.With("token", token)

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
		logger.Errorw("failed to parse uuid", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("uid", uid)

	// parse body data
	logger.Infow("parsing request body")
	var body UpdatePostHandlerRequestBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.Errorw("failed to parse request body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("body", body)

	// validate body data
	logger.Infow("validating request body")
	err = p.Validator.Struct(&body)
	if err != nil {
		logger.Errorw("failed to validate body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get post from database
	var post models.Post
	logger.Infow("getting post from database")
	err = p.MySQL.
		Model(&models.Post{}).
		Where(&models.Post{Base: models.Base{ID: uid}}).
		First(&post).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Errorw("failed to find post record in database with provided id", "err", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		logger.Errorw("failed to find post record in database", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger = logger.With("post", post)

	// check if user is post author
	logger.Infow("checking if user is author a post")
	if token.UserID != post.UserID {
		logger.Errorw("user is not author of current post", "err", err)
		err := errors.New("user is not author of current post")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update post in database
	logger.Infow("updating post in database")
	result := p.MySQL.
		Model(&post).
		Updates(&models.Post{Title: body.Title, Body: body.Body})
	if result.Error != nil || result.RowsAffected == 0 {
		if result.Error == nil {
			result.Error = errors.New("record not found")
		}
		logger.Errorw("failed to update post in database", "err", result.Error)
		http.Error(w, result.Error.Error(), http.StatusBadRequest)
		return
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := UpdatePostHandlerResponseBody{
		Post:    &post,
		Message: "successfully updated post",
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
	w.WriteHeader(http.StatusOK)

	logger.Infow("successfully updated post in database")
	w.Write(b)
}
