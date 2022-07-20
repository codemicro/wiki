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

func (e *Endpoints) Get_ListAllPages(ctx *fiber.Ctx) error {
	pages, err := e.db.GetAllPages()
	if err != nil {
		return errors.WithStack(err)
	}
	return sendNode(ctx, views.AllPagesPage(pages))
}

func (e *Endpoints) NewPage(ctx *fiber.Ctx) error {
	_, loggedIn := e.checkAuth(ctx)
	if !loggedIn {
		return redirectForSignIn(ctx)
	}

	const (
		titleKey   = "title"
		contentKey = "content"
		tagKey     = "tag"
	)

	tags, err := e.db.GetAllTags()
	if err != nil {
		return errors.WithStack(err)
	}

	switch ctx.Method() {
	case fiber.MethodGet:
		return sendNode(ctx, views.EditPagePage(views.EditPageProps{
			TitleKey:      titleKey,
			ContentKey:    contentKey,
			TagKey:        tagKey,
			Tags:          tags,
			EditPageTitle: "Create new page",
		}))
	case fiber.MethodPost:
		title := ctx.FormValue(titleKey)
		content := ctx.FormValue(contentKey)
		tag := ctx.FormValue(tagKey)

		if title == "" {
			return sendNode(ctx, views.EditPagePage(views.EditPageProps{
				TitleKey:      titleKey,
				ContentKey:    contentKey,
				TagKey:        tagKey,
				Tags:          tags,
				EditPageTitle: "Create new page",

				ContentValue:  content,
				SelectedTagID: tag,

				Problem: "Title cannot be empty",
			}))
		}

		if tag != "" {
			var found bool
			for _, t := range tags {
				if t.ID == tag {
					found = true
					break
				}
			}
			if !found {
				return sendNode(ctx, views.EditPagePage(views.EditPageProps{
					TitleKey:      titleKey,
					ContentKey:    contentKey,
					TagKey:        tagKey,
					Tags:          tags,
					EditPageTitle: "Create new page",

					TitleValue:    title,
					ContentValue:  content,
					SelectedTagID: tag,

					Problem: "Unknown tag",
				}))
			}
		}

		page := &db.Page{
			ID:      shortuuid.New(),
			Title:   title,
			Content: content,
		}

		if err := e.db.CreatePage(page); err != nil {
			return errors.WithStack(err)
		}

		if tag != "" {
			if err := e.db.AssignPageToTag(page.ID, tag); err != nil {
				return errors.WithStack(err)
			}
		}

		return ctx.Redirect(urls.Make(urls.ViewPage, page.ID))
	default:
		return errors.WithStack(fmt.Errorf("unreachable code reached: unknown method %s", ctx.Method()))
	}
}
