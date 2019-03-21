package people

import (
	"github.com/aghape-pkg/address"
	"github.com/aghape-pkg/mail"
	"github.com/aghape-pkg/phone"
	"github.com/aghape-pkg/admin"
	"github.com/aghape/db"
	"github.com/aghape/plug"
)

type Plugin struct {
	plug.EventDispatcher
	db.DBNames
	admin_plugin.AdminNames
}

func (Plugin) After() []interface{} {
	return []interface{}{&address.Plugin{}, &mail.Plugin{}, &phone.Plugin{}}
}

func (p *Plugin) OnRegister(e *plug.Options) {
	admin_plugin.Events(p).InitResources(func(e *admin_plugin.AdminEvent) {
		InitResource(e.Admin)
	})

	db.Events(p).DBOnMigrate(func(e *db.DBEvent) error {
		return e.AutoMigrate(&PeopleMedia{}, &People{}, &PeoplePhone{}, &PeopleAddress{},
			&PeopleMail{}).Error
	})
}
