package comments

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/service"
)

type Comments struct {
	ctx *service.Service
}

func (c *Comments) Setup(ctx *service.Service) {
	c.ctx = ctx

	// configure router
	c.ctx.Router.HandleFunc("/v1/comments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			c.GetCommentsHandler(w, r)
		case http.MethodPost:
			c.CreateCommentHandler(w, r)
		}
		w.WriteHeader(http.StatusNotFound)
	})

	c.ctx.Router.HandleFunc("/v1/comments/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			c.GetCommentHandler(w, r)
		case http.MethodPut:
			c.UpdateCommentHandler(w, r)
		case http.MethodDelete:
			c.DeleteCommentHandler(w, r)
		}
		w.WriteHeader(http.StatusNotFound)
	})

}
