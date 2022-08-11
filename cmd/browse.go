package cmd

import (
	"launchdata/bubble"
	"launchdata/config"

	"github.com/spf13/cobra"
)

func browseCmd() *cobra.Command {
	cmdBrowse := &cobra.Command{
		Use:   "browse",
		Short: "",
		Long:  `TODO`,
		Run: func(cmd *cobra.Command, args []string) {
			config := config.Init(cmd)
			bubble.Run(&config)
		},
	}

	return cmdBrowse
}
