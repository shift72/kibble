package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/indiereign/shift72-kibble/kibble/config"
	"github.com/indiereign/shift72-kibble/kibble/models"
	"golang.org/x/crypto/ssh/terminal"
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

			err = login(cfg)
			if err != nil {
				fmt.Println(err)
				os.Exit(-2)
			}

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

func credentials(cfg *models.Config) (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter Username for %s: ", cfg.SiteURL)
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	return strings.TrimSpace(username), strings.TrimSpace(string(bytePassword))
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

func login(cfg *models.Config) error {

	username, password := credentials(cfg)

	body := fmt.Sprintf("{\"user\":{\"email\":\"%s\",\"password\":\"%s\"}}", username, password)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/services/users/auth/sign_in", cfg.SiteURL), strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed status: %d", resp.StatusCode)
	}

	b, _ := ioutil.ReadAll(resp.Body)

	var result struct {
		APIKey string `json:"auth_token"`
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}

	fmt.Println("login successful")

	cfg.Private.APIKey = result.APIKey
	config.SavePrivateConfig(cfg)

	return nil
}
