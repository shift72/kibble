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
	"bytes"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Join - using the separator join the strings
func Join(separator string, values ...string) string {

	l := len(values)
	cleaned := make([]string, 0)
	for i := 0; i < l; i++ {
		if len(values[i]) > 0 {
			cleaned = append(cleaned, values[i])
		}
	}

	result := ""
	l = len(cleaned)
	for i := 0; i < l; i++ {
		result += cleaned[i]
		if i+1 != l {
			result += separator
		}
	}

	return result
}

// Coalesce - return the first non empty string
func Coalesce(values ...string) string {
	for i := 0; i < len(values); i++ {
		if len(values[i]) > 0 {
			return values[i]
		}
	}
	return ""
}

// ParseIntFromSlug - return the index from a slug
func ParseIntFromSlug(slug string, index int) (int, bool) {
	p := strings.Split(slug, "/")
	i, err := strconv.Atoi(p[index])
	return i, err == nil
}

// ParseIntFromString - return a int where possible
func ParseIntFromString(data string) int {
	var buffer bytes.Buffer

	for i := range data {
		r := rune(data[i])
		if unicode.IsSpace(r) {
			// skip
		} else if unicode.IsDigit(r) {
			buffer.WriteRune(r)
		} else {
			break
		}
	}

	c, err := strconv.Atoi(buffer.String())
	if err != nil {
		return 0
	}
	return c
}

// ParseTimeFromString attempts to parse the string value into a time.Time struct
func ParseTimeFromString(str string) time.Time {
	var t time.Time
	if str == "" {
		return t
	}

	t, _ = time.Parse(time.RFC3339, str)
	return t
}

// Append to list if not already present (to ensure uniqueness)
func AppendUnique(item string, list []string) []string {
	for _, x := range list {
		if x == item {
			return list
		}
	}
	return append(list, item)
}
