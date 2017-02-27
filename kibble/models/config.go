package models

// Config -
type Config struct {
	DefaultLanguage string
	Languages       map[string]string
	Routes          []Route
}

// LoadConfig - loaded the config
func LoadConfig() *Config {
	return &Config{

		DefaultLanguage: "en",
		Languages: map[string]string{
			"en": "en_US",
			"fr": "fr_FR",
		},
		Routes: []Route{
			{
				Name:         "filmIndex",
				URLPath:      "/film",
				TemplatePath: "film/index.jet",
				DataSource:   "FilmCollection",
			},
			{
				Name:         "filmItem",
				URLPath:      "/film/:filmID",
				TemplatePath: "film/item.jet",
				DataSource:   "Film",
			},
			// {
			// 	Name:         "filmItemPartial",
			// 	URLPath:      "/film/:filmID/partial.html",
			// 	TemplatePath: "film/partial.jet",
			// 	DataSource:   "Film",
			// },
		},
	}
}
