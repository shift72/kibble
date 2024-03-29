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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	assert.Equal(t, "a", Join(", ", "a"), "test 1")
	assert.Equal(t, "a", Join(", ", "a", ""), "test 2")
	assert.Equal(t, "a", Join(", ", "", "a", ""), "test 3")
	assert.Equal(t, "a, b", Join(", ", "a", "b"), "test 4")
	assert.Equal(t, "a, b", Join(", ", "a", "", "b"), "test 5")
	assert.Equal(t, "a, b", Join(", ", "", "a", "", "b"), "test 6")
	assert.Equal(t, "a, b", Join(", ", "", "a", "", "b", ""), "test 7")
}

func TestCoalesce(t *testing.T) {
	assert.Equal(t, "a", Coalesce("a"), "test 1")
	assert.Equal(t, "a", Coalesce("", "a"), "test 2")
	assert.Equal(t, "a", Coalesce("", "a", "b"), "test 3")
	assert.Equal(t, "a", Coalesce("a", "b"), "test 4")
	assert.Equal(t, "", Coalesce("", ""), "test 5")
}

func TestParseIntFromString(t *testing.T) {
	assert.Equal(t, 0, ParseIntFromString("a"), "test 1")
	assert.Equal(t, 1, ParseIntFromString("1"), "test 2")
	assert.Equal(t, 103, ParseIntFromString("103"), "test 3")
	assert.Equal(t, 103, ParseIntFromString("103.123"), "test 4")
	assert.Equal(t, 103, ParseIntFromString(" 103.123"), "test 5")
}

func TestAppendUnique(t *testing.T) {
	list := []string{"b", "c"}
	assert.Equal(t, 3, len(AppendUnique("a", list)), "Appended")
	assert.Equal(t, 2, len(AppendUnique("b", list)), "Not appended")
}

func TestParseTimeFromString(t *testing.T) {
	dates := []string{
		"2021-04-01T01:03:05Z", "2021-04-01T01:03:05.000000000+00:00", "2021-04-01T01:03:05+00:00", "2021-04-01T01:03:05",
	}

	for _, d := range dates {
		parsedDate := ParseTimeFromString(d)

		assert.Equal(t, 2021, parsedDate.Year(), d)
		assert.Equal(t, time.April, parsedDate.Month(), d)
		assert.Equal(t, 1, parsedDate.Day(), d)
	}
}
