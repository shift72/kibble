package utils

import (
	"testing"

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
