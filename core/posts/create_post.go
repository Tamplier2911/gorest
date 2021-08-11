package posts

import (
	"encoding/json"
	"encoding/xml"
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
	Post    *Post  `json:"post" xml:"post"`
	Message string `json:"message" xml:"message"`
}

// Creates post instance and stores it in database
func (p *Posts) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	logger := p.ctx.Logger.Named("CreatePostHandler")

	// parse body data
	logger.Infow("parsing request body")
	var body CreatePostRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.Errorw("failed to parse request body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("body", body)

	// parse uuid id
	logger.Infow("parsing uuid from path")
	uid, err := uuid.Parse(string(body.UserID))
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
		Title:  body.Title,
		Body:   body.Body,
	}
	err = p.ctx.MySQL.Model(&Post{}).Create(&post).Error
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
	w.WriteHeader(http.StatusCreated)

	logger.Debugw("successfully created post record in database")
	w.Write(b)
}
