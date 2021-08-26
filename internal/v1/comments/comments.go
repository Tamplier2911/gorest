package comments

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/service"
)

type Comments struct {
	*service.Service
}

func (c Comments) Setup(s *service.Service) {
	c.Service = s

	// configure router
	c.Router.HandleFunc("/api/v1/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			c.GetCommentsHandler(w, r)
		case http.MethodPost:
			c.CreateCommentHandler(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	c.Router.HandleFunc("/api/v1/comments/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			c.GetCommentHandler(w, r)
		case http.MethodPut:
			c.UpdateCommentHandler(w, r)
		case http.MethodDelete:
			c.DeleteCommentHandler(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

}
