package views

import (
	"github.com/codemicro/wiki/wiki/urls"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func IndexPage(isLoggedIn bool) g.Node {
	return BasePage(BasePageProps{
		BodyNodes: []g.Node{
			Container(
				H1(g.Text("Wiki")),
				P(g.Text("TODO: List tags")),
			),
			ControlBox(
				Ul(
					g.If(!isLoggedIn, Li(Anchor(urls.Make(urls.AuthLogin), g.Text("Log in")))),
					Li(Anchor(urls.Make("/~/list"), g.Text("List all articles"))),
					g.If(isLoggedIn, g.Group([]g.Node{
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
