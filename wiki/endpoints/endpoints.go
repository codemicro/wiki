package endpoints

import (
	"bytes"
	"crypto/rand"
	goalone "github.com/bwmarrin/go-alone"
	"github.com/codemicro/wiki/wiki/config"
	"github.com/codemicro/wiki/wiki/db"
	"github.com/codemicro/wiki/wiki/urls"
	"github.com/codemicro/wiki/wiki/util"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	"github.com/pkg/errors"
	saml "github.com/russellhaering/gosaml2"
	"time"
)

const (
	sessionCookieKey = "cdmwiki_session"
	sessionValidFor  = time.Hour * 24 * 7
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

	app.Get(urls.Index, e.Index)

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
