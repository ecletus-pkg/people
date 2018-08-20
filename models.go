package people

import (
	"strings"
	"time"

	"github.com/aghape-pkg/address"
	"github.com/aghape-pkg/mail"
	"github.com/aghape-pkg/phone"
	"github.com/aghape/db/common"
	"github.com/aghape/db/common/mixins"
	"github.com/aghape/db/common/utils"
	"github.com/aghape/fragment"
	"github.com/aghape/media/media_library"
	"github.com/aghape/validations"
	"github.com/moisespsena-go/aorm"
)

type PeopleInterface interface {
	common.WithID
	AvatarURL(styles ...string) string
	GetAvatar() *media_library.MediaBox
	IsBusiness() bool
}

const (
	ICON_BUSINESS = "/images/icon-manufacturer.png"
	ICON_MEN      = "/images/icon-men.png"
	ICON_WOMAN    = "/images/icon-woman.png"
)

type PeopleGetter interface {
	GetQorPeople() *QorPeople
}

type QorPeople struct {
	common.Model
	fragment.FragmentedModel
	mixins.RecordInfoMixin
	FullName               string `gorm:"size:255"`
	NickName               string `gorm:"size:255"`
	Business               bool
	NationalIdentification string `gorm:"size:255"`
	Male                   bool
	Birthday               time.Time
	Avatar                 media_library.MediaBox
	PhoneID                string `gorm:"size:24"`
	Phone                  phone.QorPhone
	MobileID               string `gorm:"size:24"`
	Mobile                 phone.QorPhone
	OtherPhones            []QorPeoplePhone
	MailID                 string `gorm:"size:24"`
	Mail                   mail.QorMail
	OtherMails             []QorPeopleMail
	MainAddressID          string `gorm:"size:24"`
	MainAddress            address.QorAddress
	OtherAdresses          []QorPeopleAddress
	Media                  []QorPeopleMedia `gorm:"foreignkey:QorPeopleID"`
	Notes                  string           `gorm:"type:text"`
}

func (p *QorPeople) BasicLabel() string {
	s := p.FullName
	if p.NickName != "" {
		s += " (" + p.NickName + ")"
	}
	return s
}

func (p *QorPeople) BasicIcon() string {
	return ""
}

func (p *QorPeople) BasicID() string {
	return p.ID
}

func (p *QorPeople) Stringify() string {
	return p.FullName
}

func (p *QorPeople) IsBusiness() bool {
	return p.Business
}

func (p *QorPeople) GetAvatar() *media_library.MediaBox {
	return &p.Avatar
}

func (p *QorPeople) AvatarURL(styles ...string) string {
	style := "main"
	if len(styles) > 0 && styles[0] != "" {
		style = styles[0]
	}

	if len(p.Avatar.Files) > 0 {
		return p.Avatar.URL(style)
	}

	if p.Business {
		return ICON_BUSINESS
	}

	if p.Male {
		return ICON_MEN
	} else {
		return ICON_WOMAN
	}
}

func (m *QorPeople) Clean(db *aorm.DB) {
	utils.TrimStrings(&m.FullName, &m.NickName)
}

func (p *QorPeople) Validate(db *aorm.DB) {
	if strings.TrimSpace(p.FullName) == "" {
		db.AddError(validations.Failed(p, "FullName", "Full Name is empty."))
	}
}

type QorPeoplePhone struct {
	phone.QorPhone
	QorPeopleID string `gorm:"size:24"`
}

type QorPeopleMail struct {
	mail.QorMail
	QorPeopleID string `gorm:"size:24"`
}

type QorPeopleAddress struct {
	address.QorAddress
	QorPeopleID string `gorm:"size:24"`
}

func (pa *QorPeopleAddress) GetQorAddress() *address.QorAddress {
	return &pa.QorAddress
}
