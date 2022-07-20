package endpoints

import (
	"bytes"
	"crypto/rand"
	goalone "github.com/bwmarrin/go-alone"
	"github.com/codemicro/wiki/wiki/config"
	"github.com/codemicro/wiki/wiki/db"
	"github.com/codemicro/wiki/wiki/urls"
	"github.com/codemicro/wiki/wiki/util"
	"github.com/codemicro/wiki/wiki/views"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	"github.com/pkg/errors"
	saml "github.com/russellhaering/gosaml2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"net/url"
	"time"
)

const (
	sessionCookieKey = "cdmwiki_session"
	sessionValidFor  = time.Hour * 24 * 7
)

var markdownRenderer = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

type Endpoints struct {
	db              *db.DB
	serviceProvider *saml.SAMLServiceProvider
	tokenGenerator  *goalone.Sword
}

func New(dbi *db.DB) (*Endpoints, error) {
	sp := &saml.SAMLServiceProvider{
		IdentityProviderSSOURL:      config.SAML.SSOURL,
		IdentityProviderIssuer:      config.SAML.EntityID,
		ServiceProviderIssuer:       urls.Make(urls.AuthSAML),
		AssertionConsumerServiceURL: urls.Make(urls.AuthSAMLInbound),
		SignAuthnRequests:           false, // TODO: implement this
		IDPCertificateStore:         config.SAML.IDPCertificates,
		NameIdFormat:                "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
	}

	key, err := dbi.GetSessionKey()
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			randomData := make([]byte, 64)
			_, _ = rand.Read(randomData)
			if err := dbi.StoreSessionKey(string(randomData)); err != nil {
				return nil, errors.WithStack(err)
			}
		} else {
			return nil, errors.WithStack(err)
		}
	}

	return &Endpoints{
		db:              dbi,
		serviceProvider: sp,
		tokenGenerator:  goalone.New([]byte(key), goalone.Timestamp),
	}, nil
}

func (e *Endpoints) SetupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          util.JSONErrorHandler,
		DisableStartupMessage: true,
	})

	app.Use(func(ctx *fiber.Ctx) error {
		signed := ctx.Cookies(sessionCookieKey)
		if signed == "" {
			ctx.Locals("isAuthed", false)
			return ctx.Next()
		}
		signedBytes := []byte(signed)

		val, err := e.tokenGenerator.Unsign(signedBytes)
		if err != nil {
			ctx.Locals("isAuthed", false)
			return ctx.Next()
		}

		if time.Now().UTC().Sub(e.tokenGenerator.Parse(signedBytes).Timestamp) > sessionValidFor {
			ctx.Locals("isAuthed", false)
			return ctx.Next()
		}

		ctx.Locals("isAuthed", true)
		ctx.Locals("currentUserID", string(val))
		return ctx.Next()
	})

	app.Get(urls.Index, e.Index)

	app.Get(urls.ViewPage, e.Get_ViewPage)
	app.Get(urls.Pages, e.Get_ListAllPages)
	app.Get(urls.NewPage, e.NewPage)
	app.Post(urls.NewPage, e.NewPage)
	app.Get(urls.EditPage, e.EditPage)
	app.Post(urls.EditPage, e.EditPage)

	app.Get(urls.NewTag, e.CreateTag)
	app.Post(urls.NewTag, e.CreateTag)
	app.Get(urls.ListTagPages, e.Get_ListTagPages)

	app.Get(urls.AuthLogin, e.Get_Login)
	app.Get(urls.AuthSAMLInitiate, e.Get_SAMLInitiate)
	app.Post(urls.AuthSAMLInbound, e.Post_SAMLInbound)

	return app
}

func sendNode(ctx *fiber.Ctx, node g.Node) error {
	ctx.Type("html")
	b := new(bytes.Buffer)
	_ = node.Render(b)
	return ctx.Send(b.Bytes())
}

func (e *Endpoints) checkAuth(ctx *fiber.Ctx) (string, bool) {
	if !ctx.Locals("isAuthed").(bool) {
		return "", false
	}

	return ctx.Locals("currentUserID").(string), true
}

func (e *Endpoints) makeLoginProps(ctx *fiber.Ctx) views.LogInControlListItemProps {
	_, lo := e.checkAuth(ctx)
	return views.LogInControlListItemProps{
		Ctx:        ctx,
		IsLoggedIn: lo,
	}
}

func redirectForSignIn(ctx *fiber.Ctx) error {
	return ctx.Redirect(urls.Make(urls.AuthLogin) + "?next=" + url.QueryEscape(ctx.OriginalURL()))
}
