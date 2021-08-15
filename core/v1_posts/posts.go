package posts_v1

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/service"
)

type Posts struct {
	ctx *service.Service
}

func (p *Posts) Setup(ctx *service.Service) {
	p.ctx = ctx

	// configure router
	p.ctx.Router.HandleFunc("/api/v1/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			p.GetPostsHandler(w, r)
		case http.MethodPost:
			p.CreatePostHandler(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	p.ctx.Router.HandleFunc("/api/v1/posts/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			p.GetPostHandler(w, r)
		case http.MethodPut:
			p.UpdatePostHandler(w, r)
		case http.MethodDelete:
			p.DeletePostHandler(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
