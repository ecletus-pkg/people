package people

import (
	"fmt"

	admin_tabs "github.com/ecletus-pkg/admin-tabs"
	"github.com/ecletus-pkg/mail"

	"github.com/ecletus-pkg/address"
	"github.com/ecletus-pkg/phone"
	"github.com/ecletus/admin"
	"github.com/ecletus/admin/admin_helpers"
	"github.com/ecletus/core"
	"github.com/ecletus/media/media_library"
	"github.com/ecletus/media/oss"
	"github.com/moisespsena-go/aorm"
)

const (
	SCHEME_INDIVIDUAL = "Individual"
	SCHEME_BUSINESS   = "Business"
	ResourceID        = "People"
)

var DEFAULT_SCHEMES_CATEGORIES = []string{admin_tabs.SCHEME_CATEGORY}

type Config struct {
	FieldName string
	Tabs      []*admin_tabs.Tab
}

func PrepareResource(res *admin.Resource) {
	Admin := res.GetAdmin()

	res.Meta(&admin.Meta{
		Name: "Male",
		Enabled: func(recorde interface{}, context *admin.Context, meta *admin.Meta) bool {
			if context.Type.Has(admin.SHOW) {
				return recorde.(*People).Male != nil
			}
			return true
		},
	})

	res.Meta(&admin.Meta{Name: "FullName", Required: true})

	res.Meta(&admin.Meta{Name: "DisplayTupleID", EncodedName: "ID", Valuer: func(record interface{}, context *core.Context) interface{} {
		return aorm.IdOf(record)
	}})

	res.Meta(&admin.Meta{Name: "Business", Enabled: func(record interface{}, context *admin.Context, meta *admin.Meta) bool {
		return record == nil || aorm.IdOf(record).IsZero()
	}})

	res.Meta(&admin.Meta{Name: "Male", Enabled: func(record interface{}, context *admin.Context, meta *admin.Meta) bool {
		if record == nil {
			return false
		}
		return !aorm.IdOf(record).IsZero() && !record.(*People).IsBusiness()
	}})

	res.Meta(&admin.Meta{Name: "Birthday", Type: "date", Enabled: func(record interface{}, context *admin.Context, meta *admin.Meta) bool {
		if record == nil {
			return false
		}
		return !aorm.IdOf(record).IsZero() && !record.(*People).IsBusiness()
	}})

	res.Meta(&admin.Meta{Name: "Stringify", Valuer: func(v interface{}, context *core.Context) interface{} {
		return fmt.Sprint(v)
	}})

	res.SetMeta(&admin.Meta{Name: "Avatar", Config: &media_library.MediaBoxConfig{}, Type: "image"})

	res.BasicLayout().Select(aorm.IQ("{}.id, {}.full_name, {}.nick_name, {}.male, {}.avatar, {}.business"))

	avatar := res.SetMeta(&admin.Meta{
		Name: "Avatar",
		Enabled: func(record interface{}, context *admin.Context, meta *admin.Meta) bool {
			if context.Action == "show" {
				return record.(*People).Avatar.FileSize > 0
			}

			return true
		}})

	oss.ImageMetaOnDefaultValue(avatar, func(e *admin.MetaValuerEvent) {
		if e.Recorde == nil {
			return
		}

		p := e.Recorde.(*People)

		if p.Business {
			e.Value = ICON_BUSINESS
			return
		}
		if p.Male != nil {
			if *p.Male {
				e.Value = ICON_MEN
			} else {
				e.Value = ICON_WOMAN
			}
		}
	})

	avatarURL := oss.ImageMetaURL(avatar, "Preview", oss.IMAGE_STYLE_PREVIEW)
	avatarURL.Label = "Avatar"
	res.GetMeta(admin.BASIC_META_ICON).SetValuer(avatarURL.Valuer)

	admin_helpers.FieldRichEditor(res, "Notes")

	res.RegisterScheme(SCHEME_INDIVIDUAL, &admin.SchemeConfig{
		Visible: true,
		Setup: func(s *admin.Scheme) {
			s.Categories = DEFAULT_SCHEMES_CATEGORIES
			s.DefaultFilter(&admin.DBFilter{
				Name: PKG+":individual",
				Handler: func(context *core.Context, db *aorm.DB) (*aorm.DB, error) {
					return db.Where(aorm.IQ("NOT {}.business")), nil
				},
			})
		},
	})

	res.RegisterScheme(SCHEME_BUSINESS, &admin.SchemeConfig{
		Visible: true,
		Setup: func(s *admin.Scheme) {
			s.Categories = DEFAULT_SCHEMES_CATEGORIES
			s.DefaultFilter(&admin.DBFilter{
				Name: PKG+":business",
				Handler: func(context *core.Context, db *aorm.DB) (*aorm.DB, error) {
					return db.Where(aorm.IQ("{}.business")), nil
				},
			})
		},
	})

	res.Order(aorm.IQ("{}.full_name ASC"))

	res.SortableAttrs("FullName")
	res.IndexAttrs(admin.BASIC_META_ICON, "FullName", "NickName", "Business")
	res.NewAttrs("FullName", "NickName", "Business")
	res.ShowAttrs(
		&admin.Section{
			Title: "Basic Information",
			Rows: [][]string{
				{"Avatar"},
				{"FullName", "NickName"},
				{"Business", "Male", "Birthday"},
				{"Mail"},
				{"Mobile"},
				{"Phone"},
				{"Doc"},
				{"MainAddress"},
			},
		},
		"OtherPhones",
		"OtherMails",
		&admin.Section{
			Title: "Adresses",
			Rows: [][]string{
				{"MainAddress"},
				{"OtherAdresses"},
			},
		},
		"Notes",
	)
	res.EditAttrs("ID", res.ShowAttrs())
	res.SearchAttrs("= Doc", "FullName", "NickName")
	res.CustomAttrs("display.tuple", "DisplayTupleID", "Stringify")

	checkError(phone.AddSubResource(nil, res, &Phone{}, "OtherPhones"))
	checkError(mail.AddMailSubResource(nil, res, &Mail{}, "OtherMails"))
	checkError(address.AddSubResource(nil, res, &Address{}, "OtherAdresses"))

	checkError(Admin.OnResourcesAdded(func(e *admin.ResourceEvent) error {
		var addressResource, phoneResource, mailResource = e.Resources[0], e.Resources[1], e.Resources[2]

		admin_helpers.SingleEditPairs(res,
			"MainAddress", addressResource,
			"Phone", phoneResource,
			"Mobile", phoneResource,
			"Mail", mailResource,
		)

		return nil
	}, address.ResourceID, phone.ResourceID, mail.ResourceID))

	admin.SetSelectOneConfigureCallback(res, func(cfg *admin.SelectOneConfig) {
		cfg.BottomSheetSelectedTemplateJS = admin.Select2SelectedItemTemplate{
			LabelFormat: "write(data.FullName); if (data.NickName) write(' ('+data.NickName+')');",
		}.Template()
	})
}

func InitResource(Admin *admin.Admin) *admin.Resource {
	return Admin.AddResource(&People{}, &admin.Config{
		Setup: func(res *admin.Resource) {
			PrepareResource(res)
		},
	})
}

func GetResource(Admin *admin.Admin) *admin.Resource {
	return Admin.GetResourceByID("People")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
