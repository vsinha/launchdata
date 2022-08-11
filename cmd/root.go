package cmd

import (
	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
)

// TODO Add a command for loading from file instead of doing an http request
func Root() *cobra.Command {
	litter.Config.HomePackage = "lib"
	litter.Config.HidePrivateFields = false

	rootCmd := &cobra.Command{
		Use:   "launchdata",
		Short: "Launchdata 🚀\nA tool to download and examine rocket launch data from Wikipedia",
	}
	rootCmd.PersistentFlags().Bool("dry-run", false, "Don't actually take any actions")

	rootCmd.AddCommand(cacheCmd())
	rootCmd.AddCommand(browseCmd())

	return rootCmd
}
