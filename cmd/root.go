package cmd

import (
	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
)

var dryrun bool

// TODO Add a command for loading from file instead of doing an http request
func Root() *cobra.Command {
	litter.Config.HomePackage = "lib"
	litter.Config.HidePrivateFields = false

	cmdCache := cmdCache()

	rootCmd := &cobra.Command{
		Use:   "launchdata",
		Short: "Launchdata 🚀\nA tool to download and examine rocket launch data from Wikipedia",
	}
	rootCmd.PersistentFlags().BoolVar(&dryrun, "dry-run", false, "Don't actually take any actions")

	rootCmd.AddCommand(cmdCache)

	return rootCmd
}
