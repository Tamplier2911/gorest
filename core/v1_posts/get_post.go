package posts_v1

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"strings"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Represent output data of GetPostHandler
type GetPostHandlerResponseBody struct {
	Post    *models.Post `json:"post" xml:"posts"`
	Message string       `json:"message" xml:"message"`
}

// Gets post by provided id from database, returns posts
func (p *Posts) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	logger := p.ctx.Logger.Named("GetPostHandler")

	// TODO: consider abstracting this to a middleware

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

	// parse uuid
	logger.Infow("parsing uuid from path")
	uid, err := uuid.Parse(id)
	if err != nil {
		logger.Errorw("failed to parse uuid", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("uid", uid)

	// retreive post from database
	logger.Infow("getting post from database")
	var post models.Post
	err = p.ctx.MySQL.Model(&models.Post{}).Where(&models.Post{Base: models.Base{ID: uid}}).First(&post).Error
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

	logger.Infow("successfully retrieved post by id from database")
	w.Write(b)
}
