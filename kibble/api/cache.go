package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/indiereign/shift72-kibble/kibble/models"
)

var cache = httpcache.Cache(httpcache.NewMemoryCache())

// CheckAdminCredentials - check that the admin credentials are valid
func CheckAdminCredentials(cfg *models.Config, runAsAdmin bool) {

	if cfg.Private.APIKey != "" {
		isAdmin, err := IsAdmin(cfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		if !isAdmin {
			fmt.Println("Error: api key has expired")
			os.Exit(-2)
		}
	}

	setCache(runAsAdmin)
}

func setCache(runAsAdmin bool) {

	if runAsAdmin {
		cache = diskcache.New(".kibble/cache/admin")
	} else {
		cache = diskcache.New(".kibble/cache")
	}
}

// IsAdmin - check auth token is valid
func IsAdmin(cfg *models.Config) (bool, error) {

	url := fmt.Sprintf("%s/services/users/auth/bouncer", cfg.SiteURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	if cfg.Private.APIKey != "" {
		req.Header.Add("x-auth-token", cfg.Private.APIKey)
		req.Header.Add("x-bypass-cache", "1")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode != http.StatusUnauthorized, nil
}

// Get - make an http request and read the response
func Get(cfg *models.Config, url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	if cfg.Private.APIKey != "" {
		req.Header.Add("x-auth-token", cfg.Private.APIKey)
		req.Header.Add("x-bypass-cache", "1")
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: httpcache.NewTransport(cache),
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed %s status code:%d", url, resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}
