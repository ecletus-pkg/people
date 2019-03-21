package people

import (
	"strings"
	"time"

	"github.com/ecletus/media/oss"

	"github.com/ecletus-pkg/address"
	"github.com/ecletus-pkg/mail"
	"github.com/ecletus-pkg/phone"
	"github.com/ecletus/db/common/utils"
	"github.com/ecletus/fragment"
	"github.com/ecletus/validations"
	"github.com/moisespsena-go/aorm"
)

const (
	ICON_BUSINESS = "/images/icon-manufacturer.png"
	ICON_MEN      = "/images/icon-men.png"
	ICON_WOMAN    = "/images/icon-woman.png"
)

type PeopleGetter interface {
	GetQorPeople() *People
}

type People struct {
	aorm.AuditedSDModel
	fragment.FragmentedModel
	FullName               string `gorm:"size:255"`
	NickName               string `gorm:"size:255"`
	Business               bool
	NationalIdentification string `gorm:"size:255"`
	Male                   *bool
	Birthday               time.Time
	Avatar                 oss.Image `sql:"type:text" image:"crop:false"`
	PhoneID                string    `gorm:"size:24"`
	Phone                  phone.Phone
	MobileID               string `gorm:"size:24"`
	Mobile                 phone.Phone
	OtherPhones            []PeoplePhone
	MailID                 string `gorm:"size:24"`
	Mail                   mail.Mail
	OtherMails             []PeopleMail
	MainAddressID          string `gorm:"size:24"`
	MainAddress            address.Address
	OtherAdresses          []PeopleAddress
	Media                  []PeopleMedia `gorm:"foreignkey:PeopleID"`
	Notes                  string        `gorm:"type:text"`
}

func (People) GetGormInlinePreloadFields() []string {
	return []string{"FullName", "MainAddress"}
}

func (p *People) String() string {
	s := p.FullName
	if p.NickName != "" {
		s += " (" + p.NickName + ")"
	}
	return s
}

func (p *People) Stringify() string {
	return p.FullName
}

func (p *People) IsBusiness() bool {
	return p.Business
}

func (m *People) Clean(db *aorm.DB) {
	utils.TrimStrings(&m.FullName, &m.NickName)
}

func (p *People) Validate(db *aorm.DB) {
	if strings.TrimSpace(p.FullName) == "" {
		db.AddError(validations.Failed(p, "FullName", "Full Name is empty."))
	}
}

type PeoplePhone struct {
	phone.Phone
	PeopleID string `gorm:"size:24"`
}

type PeopleMail struct {
	mail.Mail
	PeopleID string `gorm:"size:24"`
}

type PeopleAddress struct {
	address.Address
	PeopleID string `gorm:"size:24"`
}

func (pa *PeopleAddress) GetAddress() *address.Address {
	return &pa.Address
}
