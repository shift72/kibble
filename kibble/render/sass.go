package render

import (
	"fmt"
	"os"
	"path"

	libsass "github.com/wellington/go-libsass"
)

// Sass - render sass
func Sass(src, dst string) error {

	_, err := os.Stat(src)
	if err != nil {
		fmt.Println(err)
		return err
	}

	r, err := os.Open(src)
	if err != nil {
		return err
	}

	os.MkdirAll(path.Dir(dst), os.ModePerm)

	f, err := os.Create(dst)
	defer f.Close()
	if err != nil {
		return err
	}

	comp, err := libsass.New(f, r, libsass.IncludePaths([]string{"styles"}))
	if err != nil {
		return err
	}

	if err := comp.Run(); err != nil {
		return err
	}

	return nil
}
