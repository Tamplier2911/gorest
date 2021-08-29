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
			service.AuthWrapperDP(p.CreatePostHandler, p.Logger, p.Config, w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	p.Router.HandleFunc("/api/v1/posts/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			p.GetPostHandler(w, r)
		case http.MethodPut:
			service.AuthWrapperDP(p.UpdatePostHandler, p.Logger, p.Config, w, r)
		case http.MethodDelete:
			service.AuthWrapperDP(p.DeletePostHandler, p.Logger, p.Config, w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
