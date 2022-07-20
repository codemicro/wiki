package urls

import (
	"github.com/codemicro/wiki/wiki/config"
	"strings"
)

const (
	TagIDParameter  = "tagID"
	PageIDParameter = "pageID"
)

const (
	Index = "/"

	Auth      = "/auth"
	AuthLogin = Auth + "/login"
	AuthSAML  = Auth + "/saml2"

	AuthSAMLInitiate = AuthSAML + "/begin"
	AuthSAMLInbound  = AuthSAML + "/acs"

	Tags         = "/tags"
	NewTag       = Tags + "/new"
	ListTagPages = Tags + "/:" + TagIDParameter

	Pages      = "/pages"
	ViewPage   = "/~/:" + PageIDParameter
	NewPage    = Pages + "/new"
	EditPage   = Pages + "/:" + PageIDParameter + "/edit"
	DeletePage = Pages + "/:" + PageIDParameter + "/delete"
)

func Make(path string, subs ...string) string {
	sp := strings.Split(path, "/")
	var n int
	for i, item := range sp {
		if strings.HasPrefix(item, ":") {
			sp[i] = subs[n]
			n += 1
		}
	}
	return strings.TrimRight(config.HTTP.ExternalURL, "/") + strings.Join(sp, "/")
}
