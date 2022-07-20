package endpoints

import (
	"fmt"
	"github.com/codemicro/wiki/wiki/db"
	"github.com/codemicro/wiki/wiki/urls"
	"github.com/codemicro/wiki/wiki/views"
	"github.com/gofiber/fiber/v2"
	"github.com/lithammer/shortuuid/v4"
	"github.com/pkg/errors"
)

func (e *Endpoints) CreateTag(ctx *fiber.Ctx) error {
	_, ok := e.checkAuth(ctx)
	if !ok {
		return ctx.Redirect(urls.Make(urls.AuthLogin))
	}

	const tagNameKey = "tagName"

	switch ctx.Method() {
	case fiber.MethodGet:
		// return page content
		return sendNode(ctx, views.CreateTagPage(views.CreateTagPageProps{
			TagNameKey: tagNameKey,
		}))
	case fiber.MethodPost:
		// process things
		tagName := ctx.FormValue(tagNameKey)

		tag := &db.Tag{
			ID:   shortuuid.New(),
			Name: tagName,
		}

		if err := e.db.CreateTag(tag); err != nil {
			if errors.Is(err, db.ErrTagNameExists) {
				ctx.Status(fiber.StatusConflict)
				return sendNode(ctx, views.CreateTagPage(views.CreateTagPageProps{
					TagNameKey: tagNameKey,
					Problem:    "Tag name already in use",
				}))
			}
			return errors.WithStack(err)
		}

		return ctx.Redirect(urls.Make(urls.ListTagPages, tag.ID))
	default:
		return errors.WithStack(fmt.Errorf("unreachable code reached: unknown method %s", ctx.Method()))
	}
}
