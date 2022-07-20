package views

import (
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
						Li(Anchor(urls.Make(urls.CreateTag), g.Text("Create new tag"))),
					})),
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
				H1(g.Text("Create tag")),
				FormEl(
					Action(urls.Make(urls.CreateTag)),
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
		Title: "Create tag",
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
				),
			),
		},
		Title: "All pages",
	})
}
