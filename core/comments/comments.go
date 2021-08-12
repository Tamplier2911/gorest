package comments

import "github.com/Tamplier2911/gorest/pkg/service"

type Comments struct {
	ctx *service.Service
}

func (c *Comments) Setup(ctx *service.Service) {
	c.ctx = ctx
}
