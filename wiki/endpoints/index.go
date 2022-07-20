package endpoints

import (
	"github.com/codemicro/wiki/wiki/views"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func (e *Endpoints) Index(ctx *fiber.Ctx) error {
	_, loggedIn := e.checkAuth(ctx)

	tags, err := e.db.GetAllTags()
	if err != nil {
		return errors.WithStack(err)
	}

	return sendNode(ctx, views.IndexPage(views.IndexPageProps{
		IsLoggedIn: loggedIn,
		Tags:       tags,
	}))
}
