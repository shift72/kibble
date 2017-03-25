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
	"github.com/indiereign/shift72-kibble/kibble/server"
	"github.com/spf13/cobra"
)

var port int32
var watch bool
var serveAdmin bool

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves the current site",
	Long:  `Creates a local web server to test local template development.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.StartNew(port, watch, serveAdmin)
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int32VarP(&port, "port", "p", 8080, "Port to listen on")
	serveCmd.Flags().BoolP("launch", "l", false, "Launch the brower on start")
	serveCmd.Flags().BoolVar(&watch, "watch", false, "Watch for changes")
	serveCmd.Flags().BoolVar(&serveAdmin, "admin", false, "Serve using admin credentials")
}
