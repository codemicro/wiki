package endpoints

import (
	"bytes"
	"fmt"
	"github.com/codemicro/wiki/wiki/util"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
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

	rw := new(bytes.Buffer)

	fmt.Fprintf(rw, "NameID: %s\n", assertionInfo.NameID)

	fmt.Fprintf(rw, "Assertions:\n")

	for key, val := range assertionInfo.Values {
		fmt.Fprintf(rw, "  %s: %+v\n", key, val)
	}

	fmt.Fprintf(rw, "\n")

	fmt.Fprintf(rw, "Warnings:\n")
	fmt.Fprintf(rw, "%+v\n", assertionInfo.WarningInfo)

	return ctx.Send(rw.Bytes())
}
