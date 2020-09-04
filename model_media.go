package people

import (
	"encoding/json"
	"github.com/moisespsena-go/bid"
	"strings"

	"github.com/ecletus/core"
	"github.com/ecletus/media/media_library"
	"github.com/ecletus/validations"
	"github.com/moisespsena-go/aorm"
)

type Media struct {
	aorm.Model
	PeopleID     bid.BID
	Title        string
	SelectedType string
	File         media_library.MediaLibraryStorage
}

func (i *Media) Init(site *core.Site) {
	i.File.Init(site, aorm.InstanceOf(i, "File").MustFieldByName("File"))
}

func (i *Media) Validate(db *aorm.DB) {
	if strings.TrimSpace(i.Title) == "" {
		db.AddError(validations.Failed(i, "Title", "Titulo não pode ser vazio."))
	}
}

func (i *Media) SetSelectedType(typ string) {
	i.SelectedType = typ
}

func (i *Media) GetSelectedType() string {
	return i.SelectedType
}

func (i *Media) ScanMediaOptions(mediaOption media_library.MediaOption) error {
	if bytes, err := json.Marshal(mediaOption); err == nil {
		return i.File.Scan(bytes)
	} else {
		return err
	}
}

func (i *Media) GetMediaOption(ctx *core.Context) (mediaOption media_library.MediaOption) {
	mediaOption.Video = i.File.Video
	mediaOption.FileName = i.File.FileName
	mediaOption.URL = i.File.FullURL(ctx)
	mediaOption.OriginalURL = i.File.FullURL(ctx, "original")
	mediaOption.CropOptions = i.File.CropOptions
	mediaOption.Sizes = i.File.GetSizes()
	mediaOption.Description = i.File.Description
	return
}
