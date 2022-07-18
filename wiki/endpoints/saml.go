package endpoints

import (
	"github.com/codemicro/wiki/wiki/config"
	"github.com/codemicro/wiki/wiki/db"
	"github.com/codemicro/wiki/wiki/urls"
	"github.com/codemicro/wiki/wiki/util"
	"github.com/gofiber/fiber/v2"
	"github.com/lithammer/shortuuid/v4"
	"github.com/pkg/errors"
	"time"
)

func (e *Endpoints) Get_SAMLInitiate(ctx *fiber.Ctx) error {
	u, err := e.serviceProvider.BuildAuthURL("")
	if err != nil {
		return errors.WithStack(err)
	}
	return ctx.Redirect(u)
}

func (e *Endpoints) Post_SAMLInbound(ctx *fiber.Ctx) error {
	rawSAMLResponse := ctx.FormValue("SAMLResponse")

	assertionInfo, err := e.serviceProvider.RetrieveAssertionInfo(rawSAMLResponse)
	if err != nil || assertionInfo.WarningInfo.InvalidTime || assertionInfo.WarningInfo.NotInAudience {
		return util.NewRichError(fiber.StatusBadRequest, "unable to verify inbound SAML login", err)
	}

	var loginUser *db.User

	user, err := e.db.GetUserByExternalID(assertionInfo.NameID)
	if err == nil {
		loginUser = user
	} else if errors.Is(err, db.ErrNotFound) {
		loginUser = new(db.User)
		loginUser.ExternalID = assertionInfo.NameID
		if nameVal, found := assertionInfo.Values["http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name"]; found {
			loginUser.Name.String = nameVal.Values[0].Value
			loginUser.Name.Valid = true
		}
		loginUser.ID = shortuuid.New()

		if err := e.db.CreateUser(loginUser); err != nil {
			return errors.WithStack(err)
		}
	} else if err != nil {
		return errors.WithStack(err)
	}

	sessionToken := e.tokenGenerator.Sign([]byte(loginUser.ID))
	ctx.Cookie(&fiber.Cookie{
		Name:     sessionCookieKey,
		Value:    string(sessionToken),
		Expires:  time.Time{},
		Secure:   config.HTTP.SecureCookies,
		HTTPOnly: true,
	})

	return ctx.Redirect(urls.Index)
}
