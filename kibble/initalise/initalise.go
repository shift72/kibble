package initalise

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// NewSite - create a new site
func NewSite(force bool) {

	if !force {
		r, err := ioutil.ReadDir(".")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(r) > 0 {
			fmt.Println("Aborted: the current directory is not empty. Use --force to skip this check")
			os.Exit(1)
		}
	}

	// find template to clone
	fmt.Println("Searching for templates to clone...")

	results, err := GetTemplates()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Choose a template:\n")
	for i, r := range results.Items {
		fmt.Printf("%d) %s - %s owner:%s stars:%d \n", i+1, r.Name, r.Description, r.Owner.Login, r.StargazersCount)
	}

	// select a template
	reader := bufio.NewReader(os.Stdin)
	selectedID := 0
	for selectedID == 0 {
		fmt.Printf("Select a value in the range (1-%d): ", results.TotalCount)
		rawID, _ := reader.ReadString('\n')
		parsedID, err := strconv.Atoi(strings.Trim(rawID, "\n"))
		if err != nil || parsedID < 1 || parsedID > results.TotalCount {
			fmt.Println("invalid, try again")
		} else {
			selectedID = parsedID
		}
	}

	cloneURL := results.Items[selectedID-1].CloneURL
	fmt.Printf("\nCloning from %s\n", cloneURL)

	// clone does not include the git files (windows check?)
	// ideally should not include the git files
	cmd := exec.Command("git", "clone", "--depth=1", cloneURL, ".")
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cmd.Wait()

	// remove origin
	cmd = exec.Command("git", "remote", "remove", "origin")
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("\nSetup complete!\n\n")
	fmt.Printf("Next steps: \n")
	fmt.Printf(" 1. update site.json with the url of your site \n")
	fmt.Printf(" 2. initalise git and push to your repository\n")
	fmt.Printf("    `git init`\n")
	fmt.Printf("    `git remote add origin https://github.com/user/repo.git`\n")

	fmt.Printf(" 3. you can now run `kibble render --watch` to see your site\n")
}
