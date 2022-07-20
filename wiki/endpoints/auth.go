package endpoints

import (
	"github.com/codemicro/wiki/wiki/urls"
	"github.com/gofiber/fiber/v2"
)

func (e *Endpoints) Get_Login(ctx *fiber.Ctx) error {
	return ctx.Redirect(urls.AuthSAMLInitiate)
}
