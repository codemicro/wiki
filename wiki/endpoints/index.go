package endpoints

import (
	"github.com/codemicro/wiki/wiki/urls"
	"github.com/codemicro/wiki/wiki/views"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	elems "github.com/maragudk/gomponents/html"
)

func (e *Endpoints) Index(ctx *fiber.Ctx) error {
	return sendNode(ctx, views.BasePage(views.BasePageProps{
		BodyNodes: []g.Node{
			views.Container(
				elems.H1(g.Text("Wiki")),
				elems.P(g.Text("TODO: List tags")),
			),
			views.ControlBox(
				elems.Ul(
					elems.Li(views.Anchor(urls.Make(urls.AuthLogin), g.Text("Log in"))),
					elems.Li(views.Anchor(urls.Make("/~/list"), g.Text("List all articles"))),
				),
			),
		},
		Title: "Wiki",
	}))
}
