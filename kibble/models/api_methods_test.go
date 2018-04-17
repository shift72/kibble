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
