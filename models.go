package people

import (
	"strings"
	"time"

	"github.com/moisespsena-go/bid"

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

type People struct {
	aorm.AuditedSDModel
	fragment.FragmentedModel
	FullName      string `sql:"size:255"`
	NickName      string `sql:"size:255"`
	Business      bool   `sql:"not null"`
	Doc           string `sql:"size:255;unique_index:={} IS NOT NULL AND {} <> ''"`
	Male          *bool
	Birthday      *time.Time `sql:"type:date"`
	Avatar        oss.Image
	PhoneID       bid.BID
	Phone         *phone.Phone
	MobileID      bid.BID
	Mobile        *phone.Phone
	OtherPhones   []Phone `aorm:"fkc"`
	MailID        bid.BID
	Mail          *mail.Mail
	OtherMails    []Mail `aorm:"fkc"`
	MainAddressID bid.BID
	MainAddress   address.Address
	OtherAdresses []Address `aorm:"fkc"`
	Media         []Media `sql:"foreignkey:PeopleID;fkc"`
	Notes         string  `sql:"type:text"`
}

func (People) GetAormInlinePreloadFields() []string {
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

type Phone struct {
	phone.Phone
	PeopleID bid.BID
}

type Mail struct {
	mail.Mail
	PeopleID bid.BID
}

type Address struct {
	address.Address
	PeopleID bid.BID
}

func (pa *Address) GetAddress() *address.Address {
	return &pa.Address
}
