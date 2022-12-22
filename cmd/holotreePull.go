package cmd

import (
	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/htfs"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/pretty"
	"github.com/spf13/cobra"
)

var (
	remoteOrigin string
	pullRobot    string
	forcePull    bool
)

var holotreePullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Try to pull existing holotree catalog from remote source.",
	Long:  "Try to pull existing holotree catalog from remote source.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Holotree pull command lasted").Report()
		}
		_, holotreeBlueprint, err := htfs.ComposeFinalBlueprint(nil, pullRobot)
		pretty.Guard(err == nil, 1, "Blueprint calculation failed: %v", err)
		hash := htfs.BlueprintHash(holotreeBlueprint)
		tree, err := htfs.New()
		pretty.Guard(err == nil, 2, "%s", err)

		present := tree.HasBlueprint(holotreeBlueprint)
		if !present || forcePull {
			catalog := htfs.CatalogName(hash)
			err = operations.PullCatalog(remoteOrigin, catalog)
			pretty.Guard(err == nil, 3, "%s", err)
		}
		pretty.Ok()
	},
}

func init() {
	holotreeCmd.AddCommand(holotreePullCmd)
	holotreePullCmd.Flags().BoolVarP(&forcePull, "force", "", false, "Force pull check, even when blueprint is already present.")
	holotreePullCmd.Flags().StringVarP(&remoteOrigin, "origin", "o", "", "URL of remote origin to pull environment from.")
	holotreePullCmd.Flags().StringVarP(&pullRobot, "robot", "r", "robot.yaml", "Full path to 'robot.yaml' configuration file to export as catalog. <optional>")
	holotreePullCmd.MarkFlagRequired("origin")
}