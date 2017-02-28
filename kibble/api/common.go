package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

var cache = httpcache.Cache(httpcache.NewMemoryCache())

// LoadSite - load the complete site
func LoadSite(cfg *models.Config) (*models.Site, error) {

	start := time.Now()

	config, err := LoadConfig(cfg)
	if err != nil {
		return nil, err
	}

	fmt.Printf("loaded service config: %d\n", len(config))

	toggles, err := LoadFeatureToggles(cfg)
	if err != nil {
		return nil, err
	}

	fmt.Printf("loaded toggles: %d\n", len(toggles))

	stop := time.Now()
	fmt.Printf("--------------------\nLoad completed: %s\n", stop.Sub(start))

	return &models.Site{
		Config:  config,
		Toggles: toggles,
	}, nil
}

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

	cache := diskcache.New(".tmp")
	tp := httpcache.NewTransport(cache)

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tp,
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
