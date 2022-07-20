package endpoints

import (
	"github.com/codemicro/wiki/wiki/views"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func (e *Endpoints) Get_ListAllPages(ctx *fiber.Ctx) error {
	pages, err := e.db.GetAllPages()
	if err != nil {
		return errors.WithStack(err)
	}
	return sendNode(ctx, views.AllPagesPage(pages))
}
