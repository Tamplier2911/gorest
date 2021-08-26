package posts

import (
	"net/http"

	"github.com/Tamplier2911/gorest/pkg/service"
)

type Posts struct {
	*service.Service
}

func (p Posts) Setup(s *service.Service) {
	p.Service = s

	// configure router
	p.Router.HandleFunc("/api/v1/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			p.GetPostsHandler(w, r)
		case http.MethodPost:
			p.CreatePostHandler(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	p.Router.HandleFunc("/api/v1/posts/", func(w http.ResponseWriter, r *http.Request) {
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
