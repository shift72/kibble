package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)

var cache = httpcache.Cache(httpcache.NewMemoryCache())

// ConfigureDiskCache - set the cache
func ConfigureDiskCache(path string) {
	cache = diskcache.New(path)
}

// Get - make an http request and read the response
func Get(url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: httpcache.NewTransport(cache),
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed %s status code:%d", url, resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}
