package initalise

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// NewSite - create a new site
func NewSite() {

	//TODO: add check for empty directory

	// find template to clone
	fmt.Println("initalising new site...")

	results, err := GetTemplates()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("found %d results\n", results.TotalCount)
	for i, r := range results.Items {
		fmt.Printf("%d) %s - %s owner:%s stars:%d \n", i+1, r.Name, r.Description, r.Owner.Login, r.StargazersCount)
	}

	reader := bufio.NewReader(os.Stdin)
	selectedID := 0

	for selectedID == 0 {
		fmt.Printf("Choose a template, select a value in the range (1-%d): ", results.TotalCount)
		rawID, _ := reader.ReadString('\n')
		parsedID, err := strconv.Atoi(strings.Trim(rawID, "\n"))
		if err != nil {
			fmt.Println("invalid selection")
		}

		if parsedID >= 1 && parsedID <= results.TotalCount {
			selectedID = parsedID
		} else {
			fmt.Println("invalid selection")
		}
	}

	cloneURL := results.Items[selectedID-1].CloneURL
	fmt.Printf("cloning from %s\n", cloneURL)

	// git --git-dir=/dev/null clone --depth=1 /url/to/repo
	cmd := exec.Command("git", "--git-dir=/dev/null", "clone", "--depth=1", cloneURL, ".")
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

	fmt.Printf("\nYou can now run `kibble render --watch` to see your site\n")
}
