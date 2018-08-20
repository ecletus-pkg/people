package people

import (
	"github.com/aghape-pkg/address"
	"github.com/aghape-pkg/mail"
	"github.com/aghape-pkg/phone"
	"github.com/aghape/admin/adminplugin"
	"github.com/aghape/db"
	"github.com/aghape/plug"
)

type Plugin struct {
	plug.EventDispatcher
	db.DBNames
	adminplugin.AdminNames
}

func (Plugin) After() []interface{} {
	return []interface{}{&address.Plugin{}, &mail.Plugin{}, &phone.Plugin{}}
}

func (p *Plugin) OnRegister(e *plug.Options) {
	p.AdminNames.OnInitResources(p, func(e *adminplugin.AdminEvent) {
		InitResource(e.Admin)
	})

	db.Events(p).DBOnMigrateGorm(func(e *db.GormDBEvent) error {
		return e.DB.AutoMigrate(&QorPeopleMedia{}, &QorPeople{}, &QorPeoplePhone{}, &QorPeopleAddress{},
			&QorPeopleMail{}).Error
	})
}