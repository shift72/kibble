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

package render

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

// Prox - RerverseProxy object
type Prox struct {
	// target url of reverse proxy
	target *url.URL
	// instance of Go ReverseProxy thatwill do the job for us
	proxy    *httputil.ReverseProxy
	patterns []*regexp.Regexp
}

// NewProxy - create a proxy that will rewrite the scheme and host
func NewProxy(target string, patterns []string) *Prox {
	url, _ := url.Parse(target)
	director := func(req *http.Request) {
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.Host = url.Host
	}

	res := make([]*regexp.Regexp, len(patterns))
	for i := 0; i < len(patterns); i++ {
		res[i] = regexp.MustCompile(patterns[i])
	}

	return &Prox{
		target:   url,
		proxy:    &httputil.ReverseProxy{Director: director},
		patterns: res,
	}
}

// GetMiddleware - add the proxuy middleware
func (p *Prox) GetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p.shouldProxy(r.RequestURI) {
			w.Header().Add("X-GoProxy", "GoProxy")
			p.proxy.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (p *Prox) shouldProxy(requestURI string) bool {
	if strings.HasPrefix(requestURI, "/services") {
		return true
	}

	for i := 0; i < len(p.patterns); i++ {
		if p.patterns[i] != nil && p.patterns[i].MatchString(requestURI) {
			return true
		}
	}

	return false
}
