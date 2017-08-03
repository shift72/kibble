package models

import (
	"os"
	"path/filepath"
)

// SourcePath is an absolute path to where the Site source files are located.
// Can be configured in `site.json` - `"siteRootPath"`
func (cfg Config) SourcePath() string {
	if cfg.SiteRootPath == "" || cfg.SiteRootPath == "." {
		return "."
	}

	wd, _ := os.Getwd()
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

	// We shouldnt get here, maybe this method should also return an error?
	return "."
}
