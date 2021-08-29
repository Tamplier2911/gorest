package posts

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/access"
	"github.com/Tamplier2911/gorest/pkg/models"
)

// Represent input data of CreatePostHandler
type CreatePostHandlerRequestBody struct {
	Title string `json:"title" form:"title" url:"title" binding:"required" validate:"required"`
	Body  string `json:"body" form:"body" url:"body" binding:"required" validate:"required"`
} // @CreatePostResponse

// Represent output data of CreatePostHandler
type CreatePostHandlerResponseBody struct {
	Post    *models.Post `json:"post" xml:"post"`
	Message string       `json:"message" xml:"message"`
} // @CreatePostRequest

// Creates post instance and stores it in database
func (p *Posts) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	logger := p.Logger.Named("CreatePostHandler")

	// get token from context
	token := r.Context().Value("token").(*access.Token)
	logger = logger.With("token", token)

	// parse body data
	logger.Infow("parsing request body")
	var body CreatePostHandlerRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
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

	// save instance of post in database
	logger.Infow("saving post to database")
	post := models.Post{
		UserID: token.UserID,
		Title:  body.Title,
		Body:   body.Body,
	}
	err = p.MySQL.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		logger.Errorw("failed to save post in database", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := CreatePostHandlerResponseBody{
		Post:    &post,
		Message: "successfully created post",
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

	logger.Debugw("successfully created post record in database")
	w.Write(b)
}
