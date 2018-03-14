package render

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Prox - RerverseProxy object
type Prox struct {
	// target url of reverse proxy
	target *url.URL
	// instance of Go ReverseProxy thatwill do the job for us
	proxy *httputil.ReverseProxy
}

// NewProxy - create a proxy that will rewrite the scheme and host
func NewProxy(target string) *Prox {
	url, _ := url.Parse(target)
	director := func(req *http.Request) {
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.Host = url.Host
	}

	return &Prox{target: url, proxy: &httputil.ReverseProxy{Director: director}}
}

// GetMiddleware - add the proxuy middleware
func (p *Prox) GetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/services") {
			w.Header().Add("X-GoProxy", "GoProxy")
			p.proxy.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
