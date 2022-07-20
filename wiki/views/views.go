package views

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func BasePage(bodyNodes []g.Node) g.Node {
	return Doctype(
		HTML(
			Lang("en"),

			Head(
				Meta(Charset("utf8")),
				Meta(
					Name("viewport"),
					Content("width=device-width, initial-scale=1"),
				),
				TitleEl(g.Text("Wiki")),
				Link(
					Rel("stylesheet"),
					Href("/main.css"),
				),
			),

			Body(
				bodyNodes...,
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
