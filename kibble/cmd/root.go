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

package cmd

import (
	"os"

	logging "github.com/op/go-logging"
	"github.com/spf13/cobra"
)

var runAsAdmin bool
var disableCache bool
var verbose bool
var apiKey string

var log *logging.Logger

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kibble",
	Short: "SHIFT72 Front End Development tool",
	Long: `Kibble supports developers and designers to build and test templates
for the SHIFT72 Video Platform.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Errorf("%v", err)
		os.Exit(-1)
	}
}

func init() {
	log = logging.MustGetLogger("kibble")
	RootCmd.PersistentFlags().BoolVar(&runAsAdmin, "admin", false, "Render using admin credentials")
	RootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Api key to authenicate with")
	RootCmd.PersistentFlags().BoolVar(&disableCache, "disable-cache", false, "Prevent caching")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose logging")
}
