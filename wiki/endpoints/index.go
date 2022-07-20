package endpoints

import (
	"bytes"
	"github.com/codemicro/wiki/wiki/util"
	"github.com/codemicro/wiki/wiki/views"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func (e *Endpoints) Index(ctx *fiber.Ctx) error {
	tags, err := e.db.GetAllTags()
	if err != nil {
		return errors.WithStack(err)
	}

	tagFrequencies, err := e.db.GetTagFrequencies(tags)
	if err != nil {
		return errors.WithStack(err)
	}

	indexPage, err := e.db.GetPageByID("index")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := markdownRenderer.Convert([]byte(indexPage.Content), &buf); err != nil {
		return util.NewRichError(fiber.StatusInternalServerError, "Could not render markdown", err)
	}

	return sendNode(ctx, views.IndexPage(views.IndexPageProps{
		LogInControlListItemProps: e.makeLoginProps(ctx),
		RenderedContent:           buf.String(),
		Tags:                      tags,
		TagFrequencies:            tagFrequencies,
	}))
}
