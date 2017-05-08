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
