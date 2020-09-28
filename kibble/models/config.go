package models

import (
	"os"
	"path"
	"path/filepath"
)

// Config - template configuration
// NOTE: Don't use `SiteRootPath directly`, use `Config.SourcePath()` instead.
type Config struct {
	DefaultLanguage           string            `json:"defaultLanguage"`
	Languages                 map[string]string `json:"languages"`
	Routes                    []Route           `json:"routes"`
	SiteURL                   string            `json:"siteUrl"`
	BuilderVersion            string            `json:"builderVersion"`
	Version                   string            `json:"version"`
	Name                      string            `json:"Name"`
	Private                   PrivateConfig     `json:"-"`
	DisableCache              bool              `json:"-"`
	RunAsAdmin                bool              `json:"-"`
	SkipLogin                 bool              `json:"-"`
	SiteRootPath              string            `json:"siteRootPath"`
	LiveReload                LiveReloadConfig  `json:"liveReload"`
	ProxyPatterns             []string          `json:"proxy"`
	DefaultPricingCountryCode string            `json:"defaultPricingCountryCode"`
	DefaultTimeZone           string            `json:"defaultTimeZone"`
	DefaultDateFormat         string            `json:"defaultDateFormat"`
	DefaultTimeFormat         string            `json:"defaultTimeFormat"`
}

// LiveReloadConfig - configuration options for the live_reloader
type LiveReloadConfig struct {
	IgnoredPaths []string `json:"ignoredPaths"`
}

// PrivateConfig - config loaded from
type PrivateConfig struct {
	APIKey string `json:"apikey"`
}

// BuildPath returns the build path for current config
func (cfg *Config) BuildPath() string {
	if cfg.RunAsAdmin {
		return path.Join(".kibble", "build-admin")
	}
	return path.Join(".kibble", "build")
}

// FileRootPath returns the path to be used for copying
func (cfg *Config) FileRootPath() string {
	return "./" + cfg.BuildPath() + "/"
}

// SourcePath is an absolute path to where the Site source files are located.
// Can be configured in `site.json` - `"siteRootPath"`
func (cfg Config) SourcePath() string {
	wd, _ := os.Getwd()

	if cfg.SiteRootPath == "" || cfg.SiteRootPath == "." {
		return wd
	}

	src := filepath.Join(wd, cfg.SiteRootPath)

	// make sure its a directory
	fi, err := os.Stat(src)
	if err != nil {
		log.Fatalf("'%s' is not a directory.", src)
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return src
	default:
		log.Fatalf("'%s' is not a directory", src)
	}

	// We shouldn't get here, maybe this method should also return an error?
	return wd
}

// ShortCodePath is the path where the short code templates will be stored
func (cfg Config) ShortCodePath() string {
	return filepath.Join(cfg.SourcePath(), "templates/shortcodes")
}
