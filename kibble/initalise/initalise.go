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

package initalise

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"kibble/utils"
)

// NewSite - create a new site
func NewSite() {

	// set the log format for interactiveness
	utils.ConfigureInteractiveLogging(utils.ConvertToLoggingLevel(false))

	r, err := ioutil.ReadDir(".")
	if err != nil {
		log.Error("Checking directory", err)
		os.Exit(1)
	}

	if len(r) > 0 {
		log.Error("Aborted: the current directory is not empty.")
		os.Exit(1)
	}

	// find template to clone
	log.Notice("Searching for templates to clone...")

	results, err := GetTemplates()
	if err != nil {
		log.Error("Getting templates", err)
		os.Exit(1)
	}

	log.Notice("Choose a template:")
	for i, r := range results.Items {
		log.Noticef("%d) %s - %s owner:%s stars:%d", i+1, r.Name, r.Description, r.Owner.Login, r.StargazersCount)
	}

	// select a template
	reader := bufio.NewReader(os.Stdin)
	selectedID := 0
	for selectedID == 0 {
		log.Noticef("Select a value in the range (1-%d): ", results.TotalCount)
		rawID, _ := reader.ReadString('\n')
		parsedID, err := strconv.Atoi(strings.Trim(rawID, "\r\n"))
		if err != nil || parsedID < 1 || parsedID > results.TotalCount {
			log.Error("invalid, try again")
		} else {
			selectedID = parsedID
		}
	}

	cloneURL := results.Items[selectedID-1].CloneURL

	// clone repo
	logCommand("git", "clone", "--depth=1", cloneURL, ".")
	// rename origin to upstream
	logCommand("git", "remote", "rename", "origin", "upstream")

	log.Notice("\nSetup complete!!!!\n")
	log.Notice("Next steps:")
	log.Notice(" 1. npm install")
	log.Notice(" 2. npm start")
	log.Notice(" --- ")
	log.Notice(" 3. update kibble.json with the url of your site")
	log.Notice(" 4. add a remote to your repository")
	log.Notice("    `git remote add origin https://github.com/user/repo.git`")
}

func logCommand(name string, arg ...string) {
	log.Noticef("running `%s %s`", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)
	err := cmd.Start()
	if err != nil {
		log.Errorf("error `%s %s` returned %s", name, strings.Join(arg, " "), err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		log.Errorf("error `%s %s` returned %s", name, strings.Join(arg, " "), err)
		os.Exit(1)
	}
}
