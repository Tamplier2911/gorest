package comments

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

// Represent output data of DeleteCommentHandler
type DeleteCommentHandlerResponseBody struct {
	Message string `json:"message" xml:"message"`
}

// Deletes comment by provided id from database
func (c *Comments) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	logger := c.Logger.Named("DeleteCommentHandler")

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
	commentUuid, err := uuid.Parse(id)
	if err != nil {
		logger.Errorw("failed to parse uuid", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger = logger.With("commentUuid", commentUuid)

	// getting comment from database
	var comment models.Comment
	logger.Infow("getting comment from database")
	err = c.MySQL.
		Model(&models.Comment{}).
		Where(&models.Comment{Base: models.Base{ID: commentUuid}}).
		First(&comment).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Errorw("failed to find comment record in database with provided id", "err", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		logger.Errorw("failed to find comment record in database", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger = logger.With("comment", comment)

	// check if user is comment author
	logger.Infow("checking if user is author a comment")
	if token.UserID != comment.UserID {
		logger.Errorw("user is not author of current comment", "err", err)
		err := errors.New("user is not author of current comment")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// delete comment from database
	logger.Infow("deleting comment from database")
	err = c.MySQL.Delete(&comment).Error
	if err != nil {
		logger.Errorw("failed to delete comment with provided id from database", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// assemble response body
	logger.Infow("assembling response body")
	res := DeleteCommentHandlerResponseBody{
		Message: "successfully deleted comment from database",
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

	logger.Infow("successfully deleted comment from database")
	w.Write(b)
}
