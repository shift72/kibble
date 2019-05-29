package render

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldProxyServicesRoute(t *testing.T) {
	p := &Prox{}
	assert.True(t, p.shouldProxy("/Services/users/v1/auth"))
}

func TestShouldProxyFuzzyMatch(t *testing.T) {
	p := &Prox{patterns: []*regexp.Regexp{regexp.MustCompile("/custom/")}}
	assert.True(t, p.shouldProxy("/custom/hello/world"))
	assert.True(t, p.shouldProxy("/custom/"))
}

func TestShouldProxyExactMatch(t *testing.T) {
	p := &Prox{patterns: []*regexp.Regexp{regexp.MustCompile("^/custom$")}}
	assert.True(t, p.shouldProxy("/custom"))
	assert.False(t, p.shouldProxy("/custom/hello"))
}

func TestShouldProxyCaseInsensitiveMatch(t *testing.T) {
	p := &Prox{patterns: []*regexp.Regexp{regexp.MustCompile("/custom")}}
	assert.True(t, p.shouldProxy("/CUSTOM"))
	assert.True(t, p.shouldProxy("/custOm/ApI/1111"))
}
