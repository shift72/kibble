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

func (bc bonusContentV2) mapToModel2(parentSlug string, parentImages models.ImageSet) models.BonusContent {

	b := models.BonusContent{
		Slug:     fmt.Sprintf("%s/bonus/%d", parentSlug, bc.Number),
		Number:   bc.Number,
		Title:    bc.Title,
		Runtime:  models.Runtime(bc.Runtime),
		Overview: models.ApplyContentTransforms(bc.Overview),
		Images: models.ImageSet{
			Portrait:       utils.Coalesce(bc.ImageUrls.Portrait, parentImages.Portrait),
			Landscape:      utils.Coalesce(bc.ImageUrls.Landscape, parentImages.Landscape),
			Header:         utils.Coalesce(bc.ImageUrls.Header, parentImages.Header),
			Carousel:       utils.Coalesce(bc.ImageUrls.Carousel, parentImages.Carousel),
			Background:     utils.Coalesce(bc.ImageUrls.Bg, parentImages.Background),
			Classification: utils.Coalesce(bc.ImageUrls.Classification, parentImages.Classification),
		},
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
