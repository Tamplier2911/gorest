package posts

import (
	"github.com/Tamplier2911/gorest/pkg/service"
)

type Posts struct {
	ctx *service.Service
}

func (p *Posts) Setup(ctx *service.Service) {
	p.ctx = ctx
}
