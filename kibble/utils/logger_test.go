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

package utils

import (
	"testing"

	"github.com/op/go-logging"

	"github.com/stretchr/testify/assert"
)

func TestWatchedLogging(t *testing.T) {

	unique := ConfigureWatchedLogging(logging.INFO)

	log.Critical("critical1") // logged
	log.Critical("critical1") // not uniuq
	log.Error("error1")       // logged
	log.Warning("warn1")      // logged
	log.Notice("notice1")     // skipped
	log.Info("info1")         // skipped

	assert.Equal(t, 3, len(unique.Logs()), "store")
}
