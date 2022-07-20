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
	LogInControlListItemProps
	RenderedContent string
	Tags            []*db.Tag
	TagFrequencies  map[*db.Tag]int
}

func IndexPage(props IndexPageProps) g.Node {
	return BasePage(BasePageProps{
		BodyNodes: []g.Node{
			Container(
				//H1(g.Text("Wiki")),
				g.Raw(props.RenderedContent),
				H4(g.Text("Tags")),
				TagList(props.Tags, props.TagFrequencies),
			),
			ControlBox(
				Ul(
					LogInControlListItem(props.LogInControlListItemProps),
					Li(Anchor(urls.Make(urls.Pages), g.Text("List all pages"))),
					g.If(props.LogInControlListItemProps.IsLoggedIn, g.Group([]g.Node{
						Li(Anchor(urls.Make(urls.NewTag), g.Text("Create new tag"))),
						Li(Anchor(urls.Make(urls.NewPage), g.Text("Create new page"))),
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
	ActionURL string
	CancelURL string

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
					Action(props.ActionURL),
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
					Li(Anchor(props.CancelURL, g.Text("Cancel"))),
				),
			),
		},
		Title: props.EditPageTitle,
	})
}

type AllPagesPageProps struct {
	LogInControlListItemProps
	Pages []*db.Page
}

func AllPagesPage(props AllPagesPageProps) g.Node {
	sort.Slice(props.Pages, func(i, j int) bool {
		return props.Pages[i].Title < props.Pages[j].Title
	})

	return BasePage(BasePageProps{
		BodyNodes: []g.Node{
			Container(
				H1(g.Text("All pages")),
				PageTable(props.Pages),
			),
			ControlBox(
				Ul(
					LogInControlListItem(props.LogInControlListItemProps),
					g.If(
						props.LogInControlListItemProps.IsLoggedIn,
						Li(Anchor(urls.Make(urls.NewPage), g.Text("Create new page"))),
					),
					Li(Anchor(urls.Make(urls.Index), g.Text("Home"))),
				),
			),
		},
		Title: "All pages",
	})
}

type TagPagesPageProps struct {
	LogInControlListItemProps
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
					LogInControlListItem(props.LogInControlListItemProps),
					g.If(
						props.LogInControlListItemProps.IsLoggedIn,
						Li(Anchor(urls.Make(urls.NewPage), g.Text("Create new page"))),
					),
					Li(Anchor(urls.Make(urls.Index), g.Text("Home"))),
				),
			),
		},
		Title: props.Tag.Name,
	})
}

type ViewPagePageProps struct {
	LogInControlListItemProps
	Page     *db.Page
	PageTags []*db.Tag
	Rendered string
}

func ViewPagePage(props ViewPagePageProps) g.Node {
	return BasePage(BasePageProps{
		BodyNodes: []g.Node{
			Container(
				H1(g.Text(props.Page.Title)),
				P(Class("secondary"), g.Textf(
					"Created on %s, last updated on %s",
					props.Page.CreatedAt.Format("2006-01-02 15:04"),
					props.Page.UpdatedAt.Format("2006-01-02 15:04"),
				), Br(), g.Text("Tags: "), TagList(props.PageTags, nil)),
				g.Raw(props.Rendered),
			),
			ControlBox(
				Ul(
					LogInControlListItem(props.LogInControlListItemProps),
					g.If(props.LogInControlListItemProps.IsLoggedIn, g.Group([]g.Node{
						Li(Anchor(urls.Make(urls.EditPage, props.Page.ID), g.Text("Edit"))),
						Li(Anchor(urls.Make(urls.DeletePage, props.Page.ID), g.Text("Delete"))),
					})),
					Li(Anchor(urls.Make(urls.Index), g.Text("Home"))),
				),
			),
		},
		Title: props.Page.Title,
	})
}
