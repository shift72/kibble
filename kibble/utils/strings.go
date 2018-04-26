package utils

import (
	"bytes"
	"strconv"
	"strings"
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
