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

	fmt.Printf("service config: %d\n", len(config))

	toggles, err := LoadFeatureToggles(cfg)
	if err != nil {
		return nil, err
	}

	fmt.Printf("toggles: %d\n", len(toggles))

	bios, err := LoadBios(cfg)
	if err != nil {
		return nil, err
	}

	fmt.Printf("pages: %d\n", len(bios.Pages))

	//TODO: consider collecting all of the film slug
	films, err := LoadAllFilms(cfg)
	if err != nil {
		return nil, err
	}

	fmt.Printf("films: %d\n", len(films))

	stop := time.Now()
	fmt.Printf("--------------------\nLoad completed: %s\n--------------------\n", stop.Sub(start))

	return &models.Site{
		Config:     config,
		Toggles:    toggles,
		Navigation: bios.Navigation,
		Pages:      bios.Pages,
		Films:      films,
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
