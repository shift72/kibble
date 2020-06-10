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
	"errors"
	"fmt"
	"strings"
	"time"

	"kibble/config"
	"kibble/render"
	"kibble/sync"
	"kibble/utils"
	logging "github.com/op/go-logging"
	"github.com/spf13/cobra"
)

var syncCfg sync.Config
var testIdempotent bool
var renderAndSync bool

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync files to a s3 bucket",
	Long: `Syncronizes with a remove aws s3 bucket.
	Uses filename and etag to determine if the files require syncing.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// var err error
		var errCount int
		swSync := utils.NewStopwatchLevel("sync total", logging.NOTICE)

		logger := utils.ConfigureSyncLogging(utils.ConvertToLoggingLevel(verbose))

		if apiKey != "" {
			runAsAdmin = true
		}

		cfg := config.LoadConfig(runAsAdmin, apiKey, disableCache)
		config.CheckVersion(cfg)

		if testIdempotent {
			return sync.TestIdempotent(syncCfg, cfg)
		}

		var renderDuration time.Duration
		if renderAndSync {
			swRender := utils.NewStopwatchLevel("render", logging.NOTICE)
			errCount := render.Render(cfg.SourcePath(), cfg.BuildPath(), cfg)
			renderDuration = swRender.Completed()
			if errCount > 0 {
				summary := sync.Summary{}
				summary.RenderDuration = renderDuration
				summary.ErrorCount = errCount
				summary.Logs = logger.Logs()

				// return the summary to the stdout
				fmt.Println(summary.ToJSON())
				return nil
			}
		}

		// copy from the correct directory
		syncCfg.FileRootPath = cfg.FileRootPath()

		summary, err := sync.Execute(syncCfg)
		if err != nil {
			fmt.Printf("Sync failed: %s\n", err)
			return nil
		}
		summary.RenderDuration = renderDuration
		summary.ErrorCount = errCount
		summary.Logs = logger.Logs()

		// return the summary to the stdout
		fmt.Println(summary.ToJSON())

		swSync.Completed()
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if testIdempotent {
			return nil
		}

		if syncCfg.Bucket == "" {
			return errors.New("Missing argument: bucket must be set")
		}

		if syncCfg.BucketRootPath == "" {
			return errors.New("Missing argument: bucketrootpath must be set")
		}

		if !strings.HasSuffix(syncCfg.BucketRootPath, "/") {
			syncCfg.BucketRootPath = syncCfg.BucketRootPath + "/"
		}

		if strings.HasPrefix(syncCfg.BucketRootPath, "production/") {
			syncCfg.BucketRootPath = strings.Replace(syncCfg.BucketRootPath, "production/", "prod/", 1)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringVarP(&syncCfg.Profile, "profile", "p", "", "AWS Profile")
	syncCmd.Flags().StringVarP(&syncCfg.Region, "region", "r", "us-east-1", "AWS Region (default us-east-1)")
	syncCmd.Flags().StringVarP(&syncCfg.Bucket, "bucket", "b", "", "AWS S3 Bucket")
	syncCmd.Flags().StringVarP(&syncCfg.BucketRootPath, "bucketrootpath", "", "", "AWS S3 Path")

	syncCmd.Flags().BoolVarP(&renderAndSync, "render", "", false, "Renders the site before syncing.")
	syncCmd.Flags().BoolVarP(&testIdempotent, "test-idempotent", "", false, "Checks that two runs of the render process produce the same result.")
}
