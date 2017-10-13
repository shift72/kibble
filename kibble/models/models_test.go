package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildPath_As_Normal(t *testing.T) {

	cfg := &Config{
		RunAsAdmin: false,
	}
	assert.Equal(t, ".kibble/build", cfg.BuildPath())
}

func TestBuildPath_As_Admin(t *testing.T) {

	cfg := &Config{
		RunAsAdmin: true,
	}
	assert.Equal(t, ".kibble/build-admin", cfg.BuildPath())
}

func TestFileRootPath_As_Admin(t *testing.T) {

	cfg := &Config{
		RunAsAdmin: true,
	}
	assert.Equal(t, "/.kibble/build-admin/", cfg.FileRootPath())
}
