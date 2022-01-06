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

package api

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"kibble/config"
	"kibble/models"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh/terminal"
)

var cache = httpcache.Cache(httpcache.NewMemoryCache())

// CheckAdminCredentials - check that the admin credentials are valid
func CheckAdminCredentials(cfg *models.Config) {

	if cfg.RunAsAdmin && cfg.SkipLogin == false {
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

	setCache(cfg)
}

func setCache(cfg *models.Config) {
	if cfg.RunAsAdmin {
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

// returns a proxy client that is tolerant of TLS configuration
func getProxyClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
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

	resp, err := getProxyClient().Do(req)
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
	// force the cache to be bypassed regardless of being logged in or not
	req.Header.Add("x-bypass-cache", "1")
	if cfg.Private.APIKey != "" {
		req.Header.Add("x-auth-token", cfg.Private.APIKey)
	}

	client := getProxyClient()

	if !cfg.DisableCache {
		// chain the existing transport with the http cache
		cachingTransport := httpcache.NewTransport(cache)
		cachingTransport.Transport = client.Transport
		client.Transport = cachingTransport
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return nil, fmt.Errorf("request failed %s status code: %d", url, resp.StatusCode)
		}
		return nil, fmt.Errorf("request failed %s status code: %d %s", url, resp.StatusCode, string(body))
	}

	return body, err
}

// Upload a file
func Upload(cfg *models.Config, url string, params map[string]string, target string) error {
	req, err := newfileUploadRequest(url, params, "file", target)
	if err != nil {
		return err
	}

	if cfg.Private.APIKey != "" {
		req.Header.Add("x-auth-token", cfg.Private.APIKey)
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintln(body))
	}

	return nil
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	count, err := io.Copy(part, file)
	log.Debugf("uploading bytes %d", count)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func login(cfg *models.Config) error {

	username, password := credentials(cfg)

	body := fmt.Sprintf("{\"user\":{\"email\":\"%s\",\"password\":\"%s\"}}", username, password)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/services/users/auth/sign_in", cfg.SiteURL), strings.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "login failed")
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "login failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("\nlogin failed. status: %d", resp.StatusCode)
	}

	b, _ := ioutil.ReadAll(resp.Body)

	var result struct {
		APIKey string `json:"auth_token"`
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		return errors.Wrap(err, "login failed")
	}

	fmt.Println("login successful")

	cfg.Private.APIKey = result.APIKey
	config.SavePrivateConfig(cfg)

	return nil
}
