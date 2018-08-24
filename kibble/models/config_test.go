//    Copyright 2018 SHIFT72
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

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
	assert.Equal(t, "./.kibble/build-admin/", cfg.FileRootPath())
}
