package endpoints

import (
	"bytes"
	"fmt"
	"github.com/codemicro/wiki/wiki/db"
	"github.com/codemicro/wiki/wiki/urls"
	"github.com/codemicro/wiki/wiki/util"
	"github.com/codemicro/wiki/wiki/views"
	"github.com/gofiber/fiber/v2"
	"github.com/lithammer/shortuuid/v4"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
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
			ActionURL: urls.Make(urls.NewPage),
			CancelURL: urls.Make(urls.Index),
			
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
				ActionURL: urls.Make(urls.NewPage),
				CancelURL: urls.Make(urls.Index),

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
					ActionURL: urls.Make(urls.NewPage),
					CancelURL: urls.Make(urls.Index),

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

var markdownRenderer = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

func (e *Endpoints) Get_ViewPage(ctx *fiber.Ctx) error {
	pageID := ctx.Params(urls.PageIDParameter)
	page, err := e.db.GetPageByID(pageID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return fiber.ErrNotFound
		}
		return errors.WithStack(err)
	}

	tags, err := e.db.GetTagsByPageID(page.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	var buf bytes.Buffer
	if err := markdownRenderer.Convert([]byte(page.Content), &buf); err != nil {
		return util.NewRichError(fiber.StatusInternalServerError, "Could not render markdown", err)
	}

	return sendNode(ctx, views.ViewPagePage(views.ViewPagePageProps{
		Page:     page,
		PageTags: tags,
		Rendered: buf.String(),
	}))
}

func (e *Endpoints) EditPage(ctx *fiber.Ctx) error {
	_, loggedIn := e.checkAuth(ctx)
	if !loggedIn {
		return redirectForSignIn(ctx)
	}

	const (
		titleKey   = "title"
		contentKey = "content"
		tagKey     = "tag"
	)

	pageID := ctx.Params(urls.PageIDParameter)

	page, err := e.db.GetPageByID(pageID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return fiber.ErrNotFound
		}
		return errors.WithStack(err)
	}

	pageTags, err := e.db.GetTagsByPageID(pageID)
	if err != nil {
		return errors.WithStack(err)
	}
	var selectedTagID string
	if len(pageTags) != 0 {
		selectedTagID = pageTags[0].ID
	}

	tags, err := e.db.GetAllTags()
	if err != nil {
		return errors.WithStack(err)
	}

	switch ctx.Method() {
	case fiber.MethodGet:
		return sendNode(ctx, views.EditPagePage(views.EditPageProps{
			ActionURL: urls.Make(urls.EditPage, pageID),
			CancelURL: urls.Make(urls.ViewPage, pageID),

			TitleKey:      titleKey,
			ContentKey:    contentKey,
			TagKey:        tagKey,
			Tags:          tags,
			EditPageTitle: fmt.Sprintf("Edit \"%s\"", page.Title),

			TitleValue:    page.Title,
			ContentValue:  page.Content,
			SelectedTagID: selectedTagID,
		}))
	case fiber.MethodPost:
		title := ctx.FormValue(titleKey)
		content := ctx.FormValue(contentKey)
		tag := ctx.FormValue(tagKey)

		if title == "" {
			return sendNode(ctx, views.EditPagePage(views.EditPageProps{
				ActionURL: urls.Make(urls.EditPage, pageID),
				CancelURL: urls.Make(urls.ViewPage, pageID),

				TitleKey:      titleKey,
				ContentKey:    contentKey,
				TagKey:        tagKey,
				Tags:          tags,
				EditPageTitle: fmt.Sprintf("Edit \"%s\"", page.Title),

				TitleValue:    page.Title,
				ContentValue:  page.Content,
				SelectedTagID: selectedTagID,

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
					ActionURL: urls.Make(urls.EditPage, pageID),
					CancelURL: urls.Make(urls.ViewPage, pageID),

					TitleKey:      titleKey,
					ContentKey:    contentKey,
					TagKey:        tagKey,
					Tags:          tags,
					EditPageTitle: fmt.Sprintf("Edit \"%s\"", page.Title),

					TitleValue:    page.Title,
					ContentValue:  page.Content,
					SelectedTagID: selectedTagID,

					Problem: "Unknown tag",
				}))
			}
		}

		updatedPage := &db.Page{
			ID:      page.ID,
			Title:   title,
			Content: content,
		}

		if err := e.db.UpdatePage(updatedPage); err != nil {
			return errors.WithStack(err)
		}

		if tag != "" {
			if err := e.db.AssignPageToTag(updatedPage.ID, tag); err != nil {
				return errors.WithStack(err)
			}
		}

		return ctx.Redirect(urls.Make(urls.ViewPage, updatedPage.ID))
	default:
		return errors.WithStack(fmt.Errorf("unreachable code reached: unknown method %s", ctx.Method()))
	}
}
