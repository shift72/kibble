package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/indiereign/shift72-kibble/kibble/config"
	"github.com/indiereign/shift72-kibble/kibble/render"
	"github.com/indiereign/shift72-kibble/kibble/sync"
	"github.com/indiereign/shift72-kibble/kibble/utils"
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
		var err error
		swSync := utils.NewStopwatchLevel("sync total", logging.NOTICE)

		logger := utils.ConfigureSyncLogging(verbose)

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
			err := render.Render(cfg.SourcePath(), cfg.BuildPath(), cfg)
			renderDuration = swRender.Completed()
			if err != nil {
				fmt.Println("Render failed:", err)
				return err
			}
		}

		// copy from the correct directory
		syncCfg.FileRootPath = cfg.FileRootPath()

		summary, err := sync.Execute(syncCfg)
		if err != nil {
			fmt.Println("Sync failed:", err)
			return err
		}
		summary.RenderDuration = renderDuration
		summary.Errors = logger.Logs()

		// return the summary to the stdout
		fmt.Println(summary.ToJSON())

		swSync.Completed()
		return err
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
