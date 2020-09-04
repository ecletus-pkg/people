package people

import (
	"github.com/ecletus-pkg/address"
	"github.com/ecletus-pkg/mail"
	"github.com/ecletus-pkg/phone"
	"github.com/ecletus-pkg/admin"
	"github.com/ecletus/db"
	"github.com/ecletus/plug"
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
		return e.AutoMigrate(&Media{}, &People{}, &Phone{}, &Address{},
			&Mail{}).Error
	})
}
