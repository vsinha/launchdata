package cmd

import (
	"github.com/spf13/cobra"
)

func cmdBrowse() *cobra.Command {
	cmdBrowse := &cobra.Command{
		Use:   "browse",
		Short: "",
		Long:  `TODO`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return cmdBrowse
}
