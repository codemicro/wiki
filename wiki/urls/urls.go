package urls

import (
	"github.com/codemicro/wiki/wiki/config"
	"strings"
)

const (
	Index = "/"

	Auth      = "/auth"
	AuthLogin = Auth + "/login"
	AuthSAML  = Auth + "/saml2"

	AuthSAMLInitiate = AuthSAML + "/begin"
	AuthSAMLInbound  = AuthSAML + "/acs"

	Tags      = "/tags"
	CreateTag = Tags + "/new"
)

func Make(path string) string {
	return strings.TrimRight(config.HTTP.ExternalURL, "/") + path
}
