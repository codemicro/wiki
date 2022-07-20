package endpoints

import (
	"github.com/codemicro/wiki/wiki/views"
	"github.com/gofiber/fiber/v2"
)

func (e *Endpoints) Index(ctx *fiber.Ctx) error {
	_, loggedIn := e.checkAuth(ctx)
	return sendNode(ctx, views.IndexPage(loggedIn))
}
