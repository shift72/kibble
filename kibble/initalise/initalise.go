package initalise

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/indiereign/shift72-kibble/kibble/utils"
)

// NewSite - create a new site
func NewSite(force bool) {

	// set the log format for interactiveness
	utils.ConfigureInteractiveLogging(false)

	if !force {
		r, err := ioutil.ReadDir(".")
		if err != nil {
			log.Error("Checking directory", err)
			os.Exit(1)
		}

		if len(r) > 0 {
			log.Error("Aborted: the current directory is not empty. Use --force to skip this check")
			os.Exit(1)
		}
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
		parsedID, err := strconv.Atoi(strings.Trim(rawID, "\n"))
		if err != nil || parsedID < 1 || parsedID > results.TotalCount {
			log.Error("invalid, try again")
		} else {
			selectedID = parsedID
		}
	}

	cloneURL := results.Items[selectedID-1].CloneURL
	log.Notice("\nCloning from %s\n", cloneURL)

	// clone does not include the git files (windows check?)
	// ideally should not include the git files
	cmd := exec.Command("git", "clone", "--depth=1", cloneURL, ".")
	err = cmd.Start()
	if err != nil {
		log.Error("git clone failed", err)
		os.Exit(1)
	}
	cmd.Wait()

	// remove origin
	cmd = exec.Command("git", "remote", "remove", "origin")
	err = cmd.Start()
	if err != nil {
		log.Error("git clone failed", err)
		os.Exit(1)
	}

	log.Notice("Setup complete!\n")
	log.Notice("Next steps:")
	log.Notice(" 1. update site.json with the url of your site")
	log.Notice(" 2. initalise git and push to your repository")
	log.Notice("    `git init`")
	log.Notice("    `git remote add origin https://github.com/user/repo.git`")
	log.Notice(" 3. you can now run `kibble render --watch` to see your site")
}
