package main

import (
	"fmt"
	"os"

	"github.com/indiereign/shift72-kibble/kibble/cmd"
	"github.com/indiereign/shift72-kibble/kibble/datastore"
	logging "github.com/op/go-logging"
)

func init() {
	logging.SetFormatter(
		logging.MustStringFormatter(
			`%{color}%{time:15:04:05.000} ▶ %{message}%{color:reset}`,
		))
	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0))
	logging.SetLevel(logging.INFO, "")
}

func main() {

	datastore.Init()

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
