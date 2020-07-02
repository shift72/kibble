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
	"fmt"
	"sort"

	"kibble/models"
	"github.com/spf13/cobra"
)

// datasourcesCmd represents the datasource command
var datasourcesCmd = &cobra.Command{
	Use:   "datasources",
	Short: "Lists available datasources",
	Long:  `Lists available datasources and their options.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available datasources for use in routes (kibble.json)")

		ds := models.GetDataSources()

		// sort by key
		names := make([]string, 0, len(ds))
		for name := range ds {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, n := range names {
			v := ds[n]
			fmt.Printf("%s\n", v.GetName())
			for _, a := range v.GetRouteArguments() {
				fmt.Printf("  %s \t %s\n", a.Name, a.Description)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(datasourcesCmd)
}
