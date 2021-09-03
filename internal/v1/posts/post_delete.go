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
	"gorm.io/gorm/clause"
)

// Represent output data of DeletePostHandler
type DeletePostHandlerResponseBody struct {
	Message string `json:"message" xml:"message"`
} // @name DeletePostResponse

// Deletes post by provided id from database
func (p *Posts) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	logger := p.Logger.Named("DeletePostsHandler")

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

	// parse uuid
	logger.Infow("parsing uuid from path")
	uid, err := uuid.Parse(id)
	if err != nil {
		logger.Errorw("failed to parse uuid", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("uid", uid)

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

	// delete post from database
	logger.Infow("deleting post from database")
	err = p.MySQL.
		Select(clause.Associations).
		Delete(&post).
		Error
	if err != nil {
		logger.Errorw("failed to delete post record from database", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
