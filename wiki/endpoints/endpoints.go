package endpoints

import (
	"github.com/codemicro/wiki/wiki/config"
	"github.com/codemicro/wiki/wiki/db"
	"github.com/codemicro/wiki/wiki/urls"
	"github.com/codemicro/wiki/wiki/util"
	"github.com/gofiber/fiber/v2"
	saml "github.com/russellhaering/gosaml2"
)

type Endpoints struct {
	db              *db.DB
	serviceProvider *saml.SAMLServiceProvider
}

func New(dbi *db.DB) *Endpoints {
	sp := &saml.SAMLServiceProvider{
		IdentityProviderSSOURL:      config.SAML.SSOURL,
		IdentityProviderIssuer:      config.SAML.EntityID,
		ServiceProviderIssuer:       urls.Make(urls.AuthSAML),
		AssertionConsumerServiceURL: urls.Make(urls.AuthSAMLInbound),
		SignAuthnRequests:           false, // TODO: implement this
		IDPCertificateStore:         config.SAML.IDPCertificates,
		NameIdFormat:                "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
	}

	return &Endpoints{
		db:              dbi,
		serviceProvider: sp,
	}
}

func (e *Endpoints) SetupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          util.JSONErrorHandler,
		DisableStartupMessage: true,
	})

	app.Get(urls.Index, e.Index)

	app.Get(urls.AuthSAMLInitiate, e.Get_SAMLInitiate)
	app.Post(urls.AuthSAMLInbound, e.Post_SAMLInbound)

	return app
}
