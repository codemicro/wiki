package endpoints

import (
	"github.com/codemicro/wiki/wiki/urls"
	"github.com/gofiber/fiber/v2"
	"net/url"
)

func (e *Endpoints) Get_Login(ctx *fiber.Ctx) error {
	nextURL := urls.Make(urls.AuthSAMLInitiate)
	if q := ctx.Query("next"); q != "" {
		nextURL += "?next=" + url.QueryEscape(q)
	}
	return ctx.Redirect(nextURL)
}
