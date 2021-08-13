package posts

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"strings"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
)

// Represent output data of DeletePostHandler
type DeletePostHandlerResponseBody struct {
	Message string `json:"message" xml:"message"`
}

// Deletes post by provided id from database
func (p *Posts) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	logger := p.ctx.Logger.Named("DeletePostsHandler")

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

	// delete post from database
	logger.Infow("deleting post from database")
	result := p.ctx.MySQL.Model(&models.Post{}).Delete(&models.Post{Base: models.Base{ID: uid}})
	if result.Error != nil || result.RowsAffected == 0 {
		if result.Error == nil {
			result.Error = errors.New("record not found")
		}
		logger.Errorw("failed to delete post with provided id from database", "err", err)
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := DeletePostHandlerResponseBody{
		Message: "successfully deleted post from database",
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

	logger.Infow("successfully deleted post from database")
	w.Write(b)
}
