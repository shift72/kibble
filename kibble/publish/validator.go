package publish

import (
	"fmt"
	"os"
	"path"
	"strings"

	"kibble/models"
)

type descriptivePath struct {
	path    string
	purpose string
}

var requiredPaths = []descriptivePath{
	descriptivePath{
		path:    "static",
		purpose: "static files need to expect to be found in here",
	},
	descriptivePath{
		path:    "styles",
		purpose: "a styles directory style files are expected to  files need to expect to be found in here",
	},
}

var expectedPaths = []descriptivePath{
	descriptivePath{
		path:    "pin.html.jet",
		purpose: "supports logging into your via a non keyboard device (Apple TV)",
	},
	descriptivePath{
		path:    "404.html.jet",
		purpose: "a custom 404 (page not found) page",
	},
	descriptivePath{
		path:    "robots.txt.jet",
		purpose: "instruct robots on how to interact with the site",
	},
	descriptivePath{
		path:    "sitemap.xml.jet",
		purpose: "instruct search engines about the pages of interest",
	},
}

// Validate checks the site for certain files and directories
func Validate(sourcePath string, cfg *models.Config) error {
	var errors strings.Builder

	for _, p := range requiredPaths {
		if _, err := os.Stat(path.Join(sourcePath, cfg.SiteRootPath, p.path)); os.IsNotExist(err) {
			fmt.Fprintf(&errors, "\nError: missing path : %s - %s", path.Join(sourcePath, cfg.SiteRootPath, p.path), p.purpose)
		}
	}

	for _, p := range expectedPaths {
		if _, err := os.Stat(path.Join(sourcePath, cfg.SiteRootPath, p.path)); os.IsNotExist(err) {
			log.Warningf("Warning: missing path : %s - %s", path.Join(sourcePath, cfg.SiteRootPath, p.path), p.purpose)
		}
	}

	if errors.Len() > 0 {
		return fmt.Errorf("Validation failed %s", errors.String())
	}

	return nil
}
