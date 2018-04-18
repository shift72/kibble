package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringCollectionString(t *testing.T) {
	strings := StringCollection{"first", "second", "third"}
	assert.Equal(t, "first, second, third", strings.String())
}

func TestStringCollectionJoin(t *testing.T) {
	strings := StringCollection{"first", "second", "third"}
	assert.Equal(t, "first|second|third", strings.Join("|"))
}

func TestRuntimeHours(t *testing.T) {
	assert.Equal(t, 1, Runtime(100).Hours())
}

func TestRuntimeHoursLessThanOne(t *testing.T) {
	assert.Equal(t, 0, Runtime(3).Hours())
}

func TestRuntimeHoursZero(t *testing.T) {
	assert.Equal(t, 0, Runtime(0).Hours())
}

func TestRuntimeMinutesLessThanAnHour(t *testing.T) {
	assert.Equal(t, 2, Runtime(2).Minutes())
}

func TestRuntimeMinutesExactlyAnHour(t *testing.T) {
	assert.Equal(t, 0, Runtime(60).Minutes())
}

func TestRuntimeMinutesZero(t *testing.T) {
	assert.Equal(t, 0, Runtime(0).Minutes())
}
