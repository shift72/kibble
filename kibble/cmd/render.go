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
	"kibble/config"
	"kibble/render"
	"kibble/utils"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var port int32
var watch bool
var serve bool

var configOverrides map[string]string = make(map[string]string)
var toggleOverrides map[string]string = make(map[string]string)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render the entire site",
	Long: `Render templates using the available data.

Kibble is used to build and develop custom sites to run on the SHIFT72 platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		logLevel := utils.ConvertToLoggingLevel(verbose)
		utils.ConfigureStandardLogging(logLevel)

		cfg := config.LoadConfig(runAsAdmin, apiKey, disableCache)
		_ = config.CheckVersion(cfg)
		for k, v := range configOverrides {
			cfg.ConfigOverrides[k] = v
		}

		for k, v := range toggleOverrides {
			value, err := strconv.ParseBool(v)
			if err != nil {
				log.Errorf("--toggle %s: invalid boolean value %s\n", k, v)
				os.Exit(1)
			}

			cfg.ToggleOverrides[k] = value
		}

		if watch || serve {
			watchLogger := utils.ConfigureWatchedLogging(logLevel)
			render.Watch(cfg.SourcePath(), cfg.BuildPath(), cfg, port, watchLogger, watch)
		} else {
			render.Render(cfg.SourcePath(), cfg.BuildPath(), cfg)
		}
	},
}

func init() {
	RootCmd.AddCommand(renderCmd)
	renderCmd.Flags().Int32VarP(&port, "port", "p", 8080, "Port to listen on")
	renderCmd.Flags().BoolVar(&watch, "watch", false, "Watch for changes")
	renderCmd.Flags().BoolVar(&serve, "serve", false, "Serve the site, but dont watch for changes")
	renderCmd.Flags().StringToStringVar(&configOverrides, "config", make(map[string]string), "Set site configuration values for this build")
	renderCmd.Flags().StringToStringVar(&toggleOverrides, "toggle", make(map[string]string), "Set site feature toggles for this build")
}
