package views

import (
	"fmt"
	"github.com/codemicro/wiki/wiki/db"
	"github.com/codemicro/wiki/wiki/urls"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
	"sort"
)

type IndexPageProps struct {
	IsLoggedIn bool
	Tags       []*db.Tag
}

func IndexPage(props IndexPageProps) g.Node {
	return BasePage(BasePageProps{
		BodyNodes: []g.Node{
			Container(
				H1(g.Text("Wiki")),
				P(g.Text("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. In nibh mauris cursus mattis molestie a iaculis at erat. At imperdiet dui accumsan sit amet nulla facilisi. Tristique magna sit amet purus. Felis bibendum ut tristique et. Cras adipiscing enim eu turpis egestas pretium aenean pharetra magna. Ut consequat semper viverra nam. Ullamcorper sit amet risus nullam eget felis. Eget dolor morbi non arcu risus. Aenean pharetra magna ac placerat vestibulum lectus mauris ultrices eros. Rhoncus aenean vel elit scelerisque mauris pellentesque. Eu scelerisque felis imperdiet proin. Pretium fusce id velit ut. Pharetra magna ac placerat vestibulum lectus mauris ultrices eros in.")),
				H4(g.Text("Tags")),
				TagTable(props.Tags),
			),
			ControlBox(
				Ul(
					g.If(!props.IsLoggedIn, Li(Anchor(urls.Make(urls.AuthLogin), g.Text("Log in")))),
					Li(Anchor(urls.Make(urls.Pages), g.Text("List all pages"))),
					g.If(props.IsLoggedIn, g.Group([]g.Node{
						Li(Anchor(urls.Make(urls.NewTag), g.Text("Create new tag"))),
						Li(Anchor(urls.Make(urls.NewPage), g.Text("Create new page"))),
					})),
					g.Text("TODO"),
				),
			),
		},
		Title: "Wiki",
	})
}

type CreateTagPageProps struct {
	TagNameKey string
	Problem    string
}

func CreateTagPage(props CreateTagPageProps) g.Node {
	return BasePage(BasePageProps{
		BodyNodes: []g.Node{
			Container(
				H1(g.Text("Create new tag")),
				FormEl(
					Action(urls.Make(urls.NewTag)),
					Method(fiber.MethodPost),
					Input(Type("text"), Name(props.TagNameKey), Placeholder("Tag name")),
					Input(Type("submit"), Value("Submit")),
				),
				g.If(props.Problem != "", P(Class("error"), g.Text(props.Problem))),
			),
			ControlBox(
				Ul(
					Li(Anchor(urls.Make(urls.Index), g.Text("Cancel"))),
				),
			),
		},
		Title: "Create new tag",
	})
}

type EditPageProps struct {
	TitleKey, ContentKey, TagKey string
	Tags                         []*db.Tag
	EditPageTitle                string

	Problem string

	TitleValue    string
	ContentValue  string
	SelectedTagID string
}

func EditPagePage(props EditPageProps) g.Node {
	sort.Slice(props.Tags, func(i, j int) bool {
		return props.Tags[i].Name < props.Tags[j].Name
	})

	return BasePage(BasePageProps{
		BodyNodes: []g.Node{
			Container(
				H1(g.Text(props.EditPageTitle)),
				FormEl(
					Action(urls.Make(urls.NewPage)),
					Method(fiber.MethodPost),
					Input(Class("full-width"), Type("text"), Name(props.TitleKey), Placeholder("Page title"), g.If(props.TitleValue != "", Value(props.TitleValue))),
					Br(),
					Textarea(Class("full-width"), Rows("25"), Name(props.ContentKey), Placeholder("Markdown page content"), g.If(props.ContentValue != "", g.Text(props.ContentValue))),
					Br(),
					Select(
						append([]g.Node{
							Option(g.Text("(untagged)"), Value(""), g.If(props.SelectedTagID == "", Selected())), Name(props.TagKey)},
							g.Map(len(props.Tags), func(i int) g.Node {
								tag := props.Tags[i]
								return Option(Value(tag.ID), g.Text(tag.Name), g.If(props.SelectedTagID == tag.ID, Selected()))
							})...)...,
					),
					Br(),
					Input(Type("submit"), Value("Submit")),
				),
				g.If(props.Problem != "", P(Class("error"), g.Text(props.Problem))),
			),
			ControlBox(
				Ul(
					Li(Anchor(urls.Make(urls.Index), g.Text("Cancel"))),
				),
			),
		},
		Title: props.EditPageTitle,
	})
}

func AllPagesPage(pages []*db.Page) g.Node {
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Title < pages[j].Title
	})

	return BasePage(BasePageProps{
		BodyNodes: []g.Node{
			Container(
				H1(g.Text("All pages")),
				PageTable(pages),
			),
			ControlBox(
				Ul(
					Li(Anchor(urls.Make(urls.Index), g.Text("Home"))),
					g.Text("TODO"),
				),
			),
		},
		Title: "All pages",
	})
}

type TagPagesPageProps struct {
	Tag   *db.Tag
	Pages []*db.Page
}

func TagPagesPage(props TagPagesPageProps) g.Node {
	sort.Slice(props.Pages, func(i, j int) bool {
		return props.Pages[i].Title < props.Pages[j].Title
	})

	return BasePage(BasePageProps{
		BodyNodes: []g.Node{
			Container(
				H1(g.Text("Tag: "+props.Tag.Name)),
				PageTable(props.Pages),
			),
			ControlBox(
				Ul(
					Li(Anchor(urls.Make(urls.Index), g.Text("Home"))),
					g.Text("TODO"),
				),
			),
		},
		Title: props.Tag.Name,
	})
}

type ViewPagePageProps struct {
	Page     *db.Page
	Rendered string
}

func ViewPagePage(props ViewPagePageProps) g.Node {
	return BasePage(BasePageProps{
		BodyNodes: []g.Node{
			Container(
				H1(g.Text(props.Page.Title)),
				P(Class("secondary"), g.Raw(fmt.Sprintf(
					"Created on %s, last updated on %s",
					props.Page.CreatedAt.Format("2006-01-02 15:04"),
					props.Page.UpdatedAt.Format("2006-01-02 15:04"),
				))),
				g.Raw(props.Rendered),
			),
			ControlBox(
				Ul(
					Li(Anchor(urls.Make(urls.Index), g.Text("Home"))),
					g.Text("TODO"),
				),
			),
		},
		Title: props.Page.Title,
	})
}
