package people

import (
	"github.com/moisespsena-go/i18n-modular/i18nmod"
	"github.com/moisespsena-go/path-helpers"
	"github.com/moisespsena/template/html/template"
)

var ImageTag, _ = template.New("qor:db:common.people.tag.image").Parse("<img src=\"{{.}}\"></img>")

var (
	PKG        = path_helpers.GetCalledDir()
	I18N_GROUP = i18nmod.PkgToGroup(PKG)
)
