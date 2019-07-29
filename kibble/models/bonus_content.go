package models

// BonusContentCollection - all bonus content for a film or season
type BonusContentCollection []BonusContent

// BonusContent - bonus content model
type BonusContent struct {
	Slug           string
	Number         int
	Title          string
	Images         ImageSet
	SubtitleTracks []SubtitleTrack
	Runtime        Runtime
	Overview       string
	CustomFields   CustomFields
}

// GetGenericItem - returns a generic item based on the film bonus
func (bonus BonusContent) GetGenericItem() GenericItem {
	return GenericItem{
		Title:     bonus.Title,
		Slug:      bonus.Slug,
		Images:    bonus.Images,
		ItemType:  "bonus",
		InnerItem: bonus,
	}
}

// GetSubtitles - translate the SubtitleTracks list into StringCollection
func (bonus BonusContent) GetSubtitles() StringCollection {
	var result StringCollection
	for _, s := range bonus.SubtitleTracks {
		result = append(result, s.Name)
	}
	return result
}
