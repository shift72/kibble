package main

import (
	"fmt"
	"os"

	"github.com/indiereign/shift72-kibble/kibble/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
