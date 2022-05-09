package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCount(t *testing.T) {
	tests := map[string]struct {
		input  GenericItems
		sep    string
		expect map[string]int
	}{
		"empty":          {input: []GenericItem{}, expect: map[string]int{"film": 0, "tvseason": 0}},
		"single film":    {input: []GenericItem{GenericItem{ItemType: "film"}}, expect: map[string]int{"film": 1, "tvseason": 0}},
		"single season":  {input: []GenericItem{GenericItem{ItemType: "tvseason"}}, expect: map[string]int{"film": 0, "tvseason": 1}},
		"multiple items": {input: []GenericItem{GenericItem{ItemType: "film"}, GenericItem{ItemType: "tvseason"}}, expect: map[string]int{"film": 1, "tvseason": 1}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.GetCount()
			assert.Equal(t, got, tc.expect, name)
		})
	}
}
