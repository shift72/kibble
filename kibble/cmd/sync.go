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
	"errors"

	"github.com/indiereign/shift72-kibble/kibble/sync"
	"github.com/spf13/cobra"
)

var cfg = sync.Config{}

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync files to a s3 bucket",
	Long: `Syncronizes with a remove aws s3 bucket.
	Uses filename and etag to determine if the files require syncing.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return sync.Execute(cfg)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if cfg.Profile == "" {
			return errors.New("Missing argument: profile must be set")
		}

		if cfg.Bucket == "" {
			return errors.New("Missing argument: bucket must be set")
		}

		if cfg.BucketRootPath == "" {
			return errors.New("Missing argument: bucketrootpath must be set")
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringVarP(&cfg.Profile, "profile", "p", "", "AWS Profile")
	syncCmd.Flags().StringVarP(&cfg.Region, "region", "r", "us-east-1", "AWS Region (default us-east-1)")
	syncCmd.Flags().StringVarP(&cfg.Bucket, "bucket", "b", "", "AWS Profile")
	syncCmd.Flags().StringVarP(&cfg.BucketRootPath, "bucketrootpath", "", "", "AWS S3 ")
	syncCmd.Flags().StringVarP(&cfg.FileRootPath, "filerootpath", "", "./.kibble/build/", "path to upload")
}
