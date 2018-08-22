package people

import (
	"encoding/json"
	"strings"

	"github.com/aghape/core"
	"github.com/aghape/core/db"
	"github.com/aghape/media/media_library"
	"github.com/aghape/validations"
	"github.com/moisespsena-go/aorm"
)

type PeopleMedia struct {
	aorm.Model
	PeopleID     string `gorm:"size:24"`
	Title        string
	SelectedType string
	File         media_library.MediaLibraryStorage `gorm:"type:text" media_library:"url:/system/{{class}}/{{primary_key}}/{{column}}.{{extension}}"`
}

func (i *PeopleMedia) Init(site core.SiteInterface) {
	i.File.Init(site, db.FieldCache.Get(i, "File"))
}

func (i *PeopleMedia) Validate(db *aorm.DB) {
	if strings.TrimSpace(i.Title) == "" {
		db.AddError(validations.Failed(i, "Title", "Titulo não pode ser vazio."))
	}
}

func (i *PeopleMedia) SetSelectedType(typ string) {
	i.SelectedType = typ
}

func (i *PeopleMedia) GetSelectedType() string {
	return i.SelectedType
}

func (i *PeopleMedia) ScanMediaOptions(mediaOption media_library.MediaOption) error {
	if bytes, err := json.Marshal(mediaOption); err == nil {
		return i.File.Scan(bytes)
	} else {
		return err
	}
}

func (i *PeopleMedia) GetMediaOption() (mediaOption media_library.MediaOption) {
	mediaOption.Video = i.File.Video
	mediaOption.FileName = i.File.FileName
	mediaOption.URL = i.File.FullURL()
	mediaOption.OriginalURL = i.File.FullURL("original")
	mediaOption.CropOptions = i.File.CropOptions
	mediaOption.Sizes = i.File.GetSizes()
	mediaOption.Description = i.File.Description
	return
}
