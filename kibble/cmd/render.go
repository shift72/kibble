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
	"github.com/indiereign/shift72-kibble/kibble/config"
	"github.com/indiereign/shift72-kibble/kibble/render"
	"github.com/indiereign/shift72-kibble/kibble/utils"
	"github.com/spf13/cobra"
)

var port int32
var watch bool

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render the entire site",
	Long: `Render templates using the available data.

Kibble is used to build and develop custom sites to run on the SHIFT72 platform.`,
	Run: func(cmd *cobra.Command, args []string) {

		if watch {
			log := utils.ConfigureWatchedLogging(verbose)
			cfg := config.LoadConfig(runAsAdmin, apiKey, disableCache)
			config.CheckVersion(cfg)
			render.Watch(cfg.SourcePath(), buildPath, cfg, port, log)
		} else {
			utils.ConfigureStandardLogging(verbose)
			cfg := config.LoadConfig(runAsAdmin, apiKey, disableCache)
			config.CheckVersion(cfg)
			render.Render(cfg.SourcePath(), buildPath, cfg)
		}
	},
}

func init() {
	RootCmd.AddCommand(renderCmd)
	renderCmd.Flags().Int32VarP(&port, "port", "p", 8080, "Port to listen on")
	renderCmd.Flags().BoolVar(&watch, "watch", false, "Watch for changes")
}
