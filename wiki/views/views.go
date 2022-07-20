package views

import (
	"fmt"
	"github.com/codemicro/wiki/wiki/db"
	"github.com/codemicro/wiki/wiki/urls"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
	"net/url"
)

type BasePageProps struct {
	BodyNodes   []g.Node
	HeadNodes   []g.Node
	Title       string
	Description string
}

func BasePage(bp BasePageProps) g.Node {
	// TODO: Use `bp.Description`
	return Doctype(
		HTML(
			Lang("en"),

			Head(append([]g.Node{
				Meta(Charset("utf8")),
				Meta(
					Name("viewport"),
					Content("width=device-width, initial-scale=1"),
				),
				TitleEl(g.Text(bp.Title)),
				Link(
					Rel("stylesheet"),
					Href(urls.Make("/main.css")),
				)},
				bp.HeadNodes...,
			)...),

			Body(
				bp.BodyNodes...,
			),
		),
	)
}

func Container(children ...g.Node) g.Node {
	return Div(
		append([]g.Node{Class("container")}, children...)...,
	)
}

func ControlBox(children ...g.Node) g.Node {
	return Div(
		append([]g.Node{
			Class("controlBox"),
			H4(g.Text("Controls")),
		}, children...)...,
	)
}

type LogInControlListItemProps struct {
	Ctx        *fiber.Ctx
	IsLoggedIn bool
}

func LogInControlListItem(props LogInControlListItemProps) g.Node {
	return g.If(
		!props.IsLoggedIn,
		Li(
			Anchor(urls.Make(urls.AuthLogin)+"?next="+url.QueryEscape(urls.Make(props.Ctx.OriginalURL())), g.Text("Log in")),
		),
	)
}

func Anchor(url string, children ...g.Node) g.Node {
	return A(append(children, Href(url))...)
}

func TagList(tags []*db.Tag, tagFrequencies map[*db.Tag]int) g.Node {
	if len(tags) == 0 {
		return g.Text("no tags found")
	}
	var nodes []g.Node
	for i, tag := range tags {
		text := tag.Name
		if tagFrequencies != nil {
			if freq, found := tagFrequencies[tag]; found {
				text += fmt.Sprintf("(%d)", freq)
			}
		}

		node := Anchor(urls.Make(urls.ListTagPages, tag.ID), g.Text(text))

		if i != len(tags)-1 {
			// for every single tag but the last
			node = g.Group([]g.Node{node, g.Text(", ")})
		}

		nodes = append(nodes, node)
	}
	return g.Group(nodes)
}

func TagTable(tags []*db.Tag) g.Node {
	return Table(
		//THead(
		//	Th(g.Text("Tag name")),
		//	Th(),
		//),
		TBody(
			g.Map(len(tags), func(i int) g.Node {
				tag := tags[i]
				return Tr(
					Td(g.Text(tag.Name)),
					Td(Anchor(urls.Make(urls.ListTagPages, tag.ID), g.Text("[View]"))),
				)
			})...,
		),
	)
}

func PageTable(pages []*db.Page) g.Node {
	return Table(
		THead(
			Th(g.Text("Title")),
			Th(g.Text("Created")),
			Th(g.Text("Updated")),
			Th(),
		),
		TBody(
			g.Map(len(pages), func(i int) g.Node {
				page := pages[i]
				return Tr(
					Td(g.Text(page.Title)),
					Td(g.Text(page.CreatedAt.Format("2006-01-02"))),
					Td(g.Text(page.UpdatedAt.Format("2006-01-02"))),
					Td(Anchor(urls.Make(urls.ViewPage, page.ID), g.Text("[View]"))),
				)
			})...,
		),
	)
}
