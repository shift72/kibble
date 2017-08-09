// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/indiereign/shift72-kibble/kibble/config"
	"github.com/indiereign/shift72-kibble/kibble/publish"
	"github.com/indiereign/shift72-kibble/kibble/utils"
	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the current version of the template to SHIFT72",
	Long:  `Publishing will upload the current template to your site.`,
	Run: func(cmd *cobra.Command, args []string) {

		utils.ConfigureInteractiveLogging(verbose)
		// force to run as admin
		runAsAdmin = true
		cfg := config.LoadConfig(runAsAdmin, apiKey, disableCache)

		err := publish.Execute(cfg.SourcePath(), "./.kibble/dist", cfg)
		if err != nil {
			fmt.Printf("Publish failed: %v", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)
}
