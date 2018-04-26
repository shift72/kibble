package models

import "github.com/indiereign/shift72-kibble/kibble/utils"

// ServiceConfig -
type ServiceConfig map[string]string

// SelectDefaultImageType - select the default image
func (cfg ServiceConfig) SelectDefaultImageType(landscape, portrait string) string {
	if cfg["default_image_type"] == "landscape" {
		return landscape
	}
	return portrait
}

// ForceAbsoluteImagePath fixes up relative image path by prefixing the `image_root_path` to it.
func (cfg ServiceConfig) ForceAbsoluteImagePath(url string) string {
	if (len(url) > 0){
		return cfg["image_root_path"] + url
	}

	return url
}

// GetSEOTitle - get the seo title
func (cfg ServiceConfig) GetSEOTitle(seoTitle, title string) string {
	return utils.Join(" ",
		cfg["seo_title_prefix"],
		utils.Coalesce(seoTitle, title),
		cfg["seo_title_suffix"])
}

// GetSiteName - get the site name
func (cfg ServiceConfig) GetSiteName() string {
	return cfg["seo_site_name"]
}

// GetKeywords - get the keywords, appending any passed keywords
func (cfg ServiceConfig) GetKeywords(keywords string) string {
	return utils.Join(", ", cfg["seo_site_keywords"], keywords)
}
