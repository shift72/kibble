package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringCollectionString(t *testing.T) {
  strings := StringCollection{ "first", "second", "third" }
  assert.Equal(t, "first, second, third", strings.String())
}

func TestStringCollectionJoin(t *testing.T) {
  strings := StringCollection{ "first", "second", "third" }
  assert.Equal(t, "first|second|third", strings.Join("|"))
}

func TestRuntimeHours(t *testing.T) {
  var runtime Runtime = 100
  assert.Equal(t, 1, runtime.Hours())
}

func TestRuntimeHoursLessThanOne(t *testing.T) {
  var runtime Runtime = 3
  assert.Equal(t, 0, runtime.Hours())
}

func TestRuntimeHoursZero(t *testing.T) {
  var runtime Runtime = 0
  assert.Equal(t, 0, runtime.Hours())
}

func TestRuntimeMinutesLessThanAnHour(t *testing.T) {
  var runtime Runtime = 2
  assert.Equal(t, 2, runtime.Minutes())
}

func TestRuntimeMinutesExactlyAnHour(t *testing.T) {
  var runtime Runtime = 60
  assert.Equal(t, 0, runtime.Minutes())
}

func TestRuntimeMinutesZero(t *testing.T) {
  var runtime Runtime = 0
  assert.Equal(t, 0, runtime.Minutes())
}
