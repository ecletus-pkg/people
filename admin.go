package people

import (
	"fmt"

	"github.com/aghape-pkg/address"
	"github.com/aghape-pkg/admin-tabs"
	"github.com/aghape-pkg/mail"
	"github.com/aghape-pkg/phone"
	"github.com/aghape/admin"
	"github.com/aghape/admin/admincommon"
	"github.com/aghape/admin/resource_callback"
	"github.com/aghape/core"
	"github.com/aghape/core/resource"
	"github.com/aghape/core/utils"
	"github.com/aghape/db/common"
	"github.com/aghape/media"
	"github.com/aghape/media/media_library"
	"github.com/moisespsena-go/aorm"
	"github.com/moisespsena/template/html/template"
)

type Config struct {
	FieldName string
	Tabs      admin_tabs.Tabs
}

var PeopleCallbacks = resource_callback.NewCallbacksStack()

func AddSubResource(res *admin.Resource, value interface{}, config ...*Config) *admin.Resource {
	if len(config) == 0 {
		config = []*Config{{}}
	}

	cfg := config[0]

	r := res.NewResource(&admin.SubConfig{FieldName: cfg.FieldName}, value)
	r.SetI18nModel(&People{})
	PrepareResource(r, cfg.Tabs)
	res.SetMeta(&admin.Meta{Name: cfg.FieldName, Resource: r})
	return r
}

func PrepareResource(res *admin.Resource, pageTabs admin_tabs.Tabs) {
	Admin := res.GetAdmin()

	admin_tabs.PrepareResource(res, pageTabs, DefaultTab)
	admincommon.RecordInfoFields(res)
	phone.AddSubResource(res, &PeoplePhone{}, "OtherPhones")
	mail.AddMailSubResource(res, &PeopleMail{}, "OtherMails")
	address.AddSubResource(res, &PeopleAddress{}, "OtherAdresses")

	addressResource := address.GetResource(Admin)
	phoneResource := phone.GetResource(Admin)
	mailResource := mail.GetResource(Admin)

	res.SetMeta(&admin.Meta{Name: "MainAddress", Type: "single_edit", Resource: addressResource})
	res.SetMeta(&admin.Meta{Name: "Phone", Type: "single_edit", Resource: phoneResource})
	res.SetMeta(&admin.Meta{Name: "Mobile", Type: "single_edit", Resource: phoneResource})
	res.SetMeta(&admin.Meta{Name: "Mail", Type: "single_edit", Resource: mailResource})
	res.SetMeta(&admin.Meta{Name: "Avatar", Config: &media_library.MediaBoxConfig{}, Type: "media_library_image"})

	res.Filter(&admin.Filter{
		Name:  "business",
		Label: "Business/Individual",
		Available: func(context *core.Context) bool {
			return admin_tabs.GetTabPath(context) == ""
		},
		Config: &admin.SelectOneConfig{
			Collection: utils.TuplesIndex("Individual", "Business"),
		},
		Handler: func(db *aorm.DB, argument *admin.FilterArgument) *aorm.DB {
			v := argument.Value.Get("Value").Value.([]string)[0]
			if v == "0" {
				return db.Where("NOT business")
			}
			return db.Where("business")
		}})

	res.GetLayout(resource.BASIC_LAYOUT).(*resource.Layout).PrepareFunc = func(crud *resource.CRUD) *resource.CRUD {
		crud.Context().DB = crud.Context().DB.Select("id, full_name, nick_name")
		return crud
	}

	mediaResource := res.AddResource(&admin.SubConfig{FieldName: "Media"}, nil, &admin.Config{Priority: -1})
	mediaResource.Filter(&admin.Filter{
		Name:       "SelectedType",
		Label:      "Media Type",
		Operations: []string{"contains"},
		Config:     &admin.SelectOneConfig{Collection: [][]string{{"video", "Video"}, {"image", "Image"}, {"file", "File"}, {"video_link", "Video Link"}}},
	})
	mediaResource.IndexAttrs("File", "Title")

	res.SetMeta(&admin.Meta{Name: "AvatarDisplayURL", Valuer: func(i interface{}, context *core.Context) interface{} {
		avatar := res.GetDefinedMeta(admin.BASIC_META_ICON).Valuer(i, context).(string)
		if avatar == "" {
			avatar = context.GenGlobalStaticURL(i.(*People).AvatarURL())
		}
		return avatar
	}})

	res.Meta(&admin.Meta{Name: "AvatarImageTag", Label: "Avatar", Valuer: func(record interface{}, context *core.Context) interface{} {
		if record != nil {
			uri := res.GetDefinedMeta("AvatarDisplayURL").Valuer(record, context).(string)
			tag, err := ImageTag.ExecuteString(uri)
			if err != nil {
				context.AddError(err)
			}
			return template.HTML(tag)
		}
		return ""
	}})

	res.SetMeta(&admin.Meta{
		Name: "Avatar",
		Enabled: func(record interface{}, context *admin.Context, meta *admin.Meta) bool {
			if context.Action == "show" {
				return len(record.(PeopleInterface).GetAvatar().Files) > 0
			}

			return true
		},
		Config: &media_library.MediaBoxConfig{
			RemoteDataResource: admin.NewDataResource(mediaResource),
			Max:                1,
			Sizes: map[string]*media.Size{
				"main": {Width: 560, Height: 700},
			},
		}})

	res.Meta(&admin.Meta{Name: "Notes", Config: &admin.RichEditorConfig{}})
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
		r := record.(PeopleInterface)
		return r.GetID() != "" && !r.IsBusiness()
	}})

	res.Meta(&admin.Meta{Name: "Birthday", Type: "date", Enabled: func(record interface{}, context *admin.Context, meta *admin.Meta) bool {
		if record == nil {
			return false
		}
		r := record.(PeopleInterface)
		return r.GetID() != "" && !r.IsBusiness()
	}})

	res.IndexAttrs("AvatarImageTag", "NickName", "FullName")
	res.EditAttrs("ID", res.ShowAttrs())
	res.NewAttrs("FullName", "NickName", "Business")
	res.SearchAttrs("FullName", "NickName")
	res.CustomAttrs("display.tuple", "DisplayTupleID", "Stringify")
	//res.MetaAliases

	for _, f := range res.FakeScope.Fields() {
		fieldName := f.Name
		res.GetMeta(fieldName).I18nPrefix = I18N_GROUP
	}

	PeopleCallbacks.Run(res)
}

func InitResource(Admin *admin.Admin) *admin.Resource {
	return Admin.AddResource(&People{}, &admin.Config{
		Setup: func(res *admin.Resource) {
			PrepareResource(res, PeopleTabs)
		},
	})
}

func GetResource(Admin *admin.Admin) *admin.Resource {
	return Admin.GetResourceByID("People")
}
