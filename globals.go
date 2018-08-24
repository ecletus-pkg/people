package people

import (
	"github.com/aghape-pkg/admin-tabs"
	"github.com/aghape/admin"
	"github.com/aghape/core"
	"github.com/aghape/core/utils"
	"github.com/moisespsena-go/aorm"
	"github.com/moisespsena/go-i18n-modular/i18nmod"
	"github.com/moisespsena/go-path-helpers"
	"github.com/moisespsena/template/html/template"
)

var DefaultTab = &admin_tabs.Tab{
	Default: true,
	Title:   "All",
	Handler: func(res *admin.Resource, context *core.Context, db *aorm.DB) *aorm.DB {
		return db
	},
}

var ImageTag, _ = template.New("qor:db:common.people.tag.image").Parse("<img src=\"{{.}}\"></img>")
var PeopleTabs = []*admin_tabs.Tab{
	DefaultTab,
	{"Individual", "", "", func(res *admin.Resource, context *core.Context, db *aorm.DB) *aorm.DB {
		return db
	}, false},
	{"Business", "", "", func(res *admin.Resource, context *core.Context, db *aorm.DB) *aorm.DB {
		return db
	}, false},
}

var (
	PKG        = path_helpers.GetCalledDir()
	I18N_GROUP = i18nmod.PkgToGroup(PKG)
)

func init() {
	group := I18N_GROUP + ".People"

	for _, scope := range PeopleTabs {
		scope.Path = utils.ToParamString(scope.Title)
		scope.TitleKey = group + ".scopes." + scope.Title
	}
}
