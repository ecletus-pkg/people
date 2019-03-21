package people

import (
	"fmt"

	"github.com/aghape/admin/admin_helpers"

	"github.com/aghape/media/oss"

	"github.com/aghape-pkg/address"
	"github.com/aghape-pkg/admin-tabs"
	"github.com/aghape-pkg/mail"
	"github.com/aghape-pkg/phone"
	"github.com/aghape/admin"
	"github.com/aghape/admin/resource_callback"
	"github.com/aghape/core"
	"github.com/aghape/db/common"
	"github.com/aghape/media/media_library"
	"github.com/moisespsena-go/aorm"
)

const (
	SCHEME_INDIVIDUAL = "Individual"
	SCHEME_BUSINESS   = "Business"
)

var DEFAULT_SCHEMES_CATEGORIES = []string{admin_tabs.SCHEME_CATEGORY}

type Config struct {
	FieldName string
	Tabs      []*admin_tabs.Tab
}

var PeopleCallbacks = resource_callback.NewCallbacksStack()

func PrepareResource(res *admin.Resource) {
	Admin := res.GetAdmin()

	//admin_tabs.PrepareResource(res, pageTabs, DefaultTab)
	phone.AddSubResource(res, &PeoplePhone{}, "OtherPhones")
	mail.AddMailSubResource(res, &PeopleMail{}, "OtherMails")
	address.AddSubResource(res, &PeopleAddress{}, "OtherAdresses")

	addressResource := address.GetResource(Admin)
	phoneResource := phone.GetResource(Admin)
	mailResource := mail.GetResource(Admin)

	res.RegisterScheme(SCHEME_INDIVIDUAL, &admin.SchemeConfig{
		Visible: true,
		Setup: func(s *admin.Scheme) {
			s.Categories = DEFAULT_SCHEMES_CATEGORIES
			s.DefaultFilter(func(context *core.Context, db *aorm.DB) *aorm.DB {
				return db.Where("NOT peoples.business")
			})
		},
	})
	res.RegisterScheme(SCHEME_BUSINESS, &admin.SchemeConfig{
		Visible: true,
		Setup: func(s *admin.Scheme) {
			s.Categories = DEFAULT_SCHEMES_CATEGORIES
			s.DefaultFilter(func(context *core.Context, db *aorm.DB) *aorm.DB {
				return db.Where("peoples.business")
			})
		},
	})

	res.Meta(&admin.Meta{Name: "Male", Enabled: func(recorde interface{}, context *admin.Context, meta *admin.Meta) bool {
		if context.Type.Has(admin.SHOW) {
			return recorde.(*People).Male != nil
		}
		return true
	}})

	admin_helpers.SingleEditPairs(res,
		"MainAddress", addressResource,
		"Phone", phoneResource,
		"Mobile", phoneResource,
		"Mail", mailResource,
	)
	
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

	res.Meta(&admin.Meta{Name: "Stringify", Valuer: func(v interface{}, context *core.Context) interface{} {
		return fmt.Sprint(v)
	}})

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
				{"NationalIdentification"},
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

	res.Meta(&admin.Meta{Name: "DisplayTupleID", EncodedName: "ID", Valuer: func(instance interface{}, context *core.Context) interface{} {
		return instance.(common.WithID).GetID()
	}})

	res.Meta(&admin.Meta{Name: "Business", Enabled: func(record interface{}, context *admin.Context, meta *admin.Meta) bool {
		return record == nil || record.(common.WithID).GetID() == ""
	}})

	res.Meta(&admin.Meta{Name: "Male", Enabled: func(record interface{}, context *admin.Context, meta *admin.Meta) bool {
		if record == nil {
			return false
		}
		r := record.(*People)
		return r.GetID() != "" && !r.IsBusiness()
	}})

	res.Meta(&admin.Meta{Name: "Birthday", Type: "date", Enabled: func(record interface{}, context *admin.Context, meta *admin.Meta) bool {
		if record == nil {
			return false
		}
		r := record.(*People)
		return r.GetID() != "" && !r.IsBusiness()
	}})

	res.Order(aorm.IQ("{}.full_name ASC"))

	res.SortableAttrs("FullName")
	res.IndexAttrs(admin.BASIC_META_ICON, "FullName", "NickName", "Business")
	res.EditAttrs("ID", res.ShowAttrs())
	res.NewAttrs("FullName", "NickName", "Business")
	res.SearchAttrs("FullName", "NickName")
	res.CustomAttrs("display.tuple", "DisplayTupleID", "Stringify")

	PeopleCallbacks.Run(res)
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
