package endpoints

import (
	"bytes"
	"github.com/codemicro/wiki/wiki/views"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	elems "github.com/maragudk/gomponents/html"
)

func (e *Endpoints) Index(ctx *fiber.Ctx) error {
	ctx.Type("html")
	b := new(bytes.Buffer)
	node := views.BasePage([]g.Node{
		views.Container(
			elems.H1(g.Text("Wiki")),
			elems.P(g.Text("TODO: List tags")),
		),
		views.ControlBox(
			elems.Ul(
				elems.Li(g.Text("Log in")),
				elems.Li(g.Text("Sitemap")),
			),
		),
	})
	_ = node.Render(b)
	return ctx.Send(b.Bytes())
}
