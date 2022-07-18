package urls

import (
	"github.com/codemicro/wiki/wiki/config"
	"strings"
)

const (
	Index = "/"

	Auth     = "/auth"
	AuthSAML = Auth + "/saml2"

	AuthSAMLInitiate = AuthSAML + "/begin"
	AuthSAMLInbound  = AuthSAML + "/acs"
)

func Make(path string) string {
	return strings.TrimRight(config.HTTP.ExternalURL, "/") + path
}
