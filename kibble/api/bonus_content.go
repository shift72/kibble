package api

import (
	"fmt"

	"kibble/models"
	"kibble/utils"
)

// bonusContentV2 - bonus content model
type bonusContentV2 struct {
	Number    int     `json:"number"`
	Title     string  `json:"title"`
	Overview  string  `json:"description"`
	Runtime   float64 `json:"runtime"`
	ImageUrls struct {
		Portrait       string `json:"portrait"`
		Landscape      string `json:"landscape"`
		Header         string `json:"header"`
		Carousel       string `json:"carousel"`
		Bg             string `json:"bg"`
		Classification string `json:"classification"`
	} `json:"image_urls"`
	SubtitleTracks []subtitleTrackV1      `json:"subtitle_tracks"`
	CustomFields   map[string]interface{} `json:"custom"`
}

func (bcv2 bonusContentV2) mapToModel2(parentSlug string, parentImages models.ImageSet) models.BonusContent {

	b := models.BonusContent{
		Slug:     fmt.Sprintf("%s/bonus/%d", parentSlug, bcv2.Number),
		Number:   bcv2.Number,
		Title:    bcv2.Title,
		Runtime:  models.Runtime(bcv2.Runtime),
		Overview: bcv2.Overview,
		Images: models.ImageSet{
			Portrait:       utils.Coalesce(bcv2.ImageUrls.Portrait, parentImages.Portrait),
			Landscape:      utils.Coalesce(bcv2.ImageUrls.Landscape, parentImages.Landscape),
			Header:         utils.Coalesce(bcv2.ImageUrls.Header, parentImages.Header),
			Carousel:       utils.Coalesce(bcv2.ImageUrls.Carousel, parentImages.Carousel),
			Background:     utils.Coalesce(bcv2.ImageUrls.Bg, parentImages.Background),
			Classification: utils.Coalesce(bcv2.ImageUrls.Classification, parentImages.Classification),
		},
		CustomFields: bcv2.CustomFields,
	}

	for _, t := range bcv2.SubtitleTracks {
		b.SubtitleTracks = append(b.SubtitleTracks, models.SubtitleTrack{
			Language: t.Language,
			Name:     t.Name,
			Type:     t.Type,
			Path:     t.Path,
		})
	}

	return b

}

// For parent models which use the ImageMap rather than ImageSet
func (bc bonusContentV2) mapToModel3(parentSlug string, parentImages models.ImageMap) models.BonusContent {

	// Convert the ImageMap to an ImageSet (which we will eventually phase out)
	images := models.ImageMapToImageSet(parentImages)
	images.Portrait = utils.Coalesce(bc.ImageUrls.Portrait, images.Portrait)
	images.Landscape = utils.Coalesce(bc.ImageUrls.Landscape, images.Landscape)
	images.Header = utils.Coalesce(bc.ImageUrls.Header, images.Header)
	images.Carousel = utils.Coalesce(bc.ImageUrls.Carousel, images.Carousel)
	images.Background = utils.Coalesce(bc.ImageUrls.Bg, images.Background)
	images.Classification = utils.Coalesce(bc.ImageUrls.Classification, images.Classification)

	b := models.BonusContent{
		Slug:         fmt.Sprintf("%s/bonus/%d", parentSlug, bc.Number),
		Number:       bc.Number,
		Title:        bc.Title,
		Runtime:      models.Runtime(bc.Runtime),
		Overview:     bc.Overview,
		Images:       images,
		CustomFields: bc.CustomFields,
	}

	for _, t := range bc.SubtitleTracks {
		b.SubtitleTracks = append(b.SubtitleTracks, models.SubtitleTrack{
			Language: t.Language,
			Name:     t.Name,
			Type:     t.Type,
			Path:     t.Path,
		})
	}

	return b

}
