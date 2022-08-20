package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"launchdata/bubble"
	"launchdata/config"

	"github.com/spf13/cobra"
)

func browseCmd() *cobra.Command {
	cmdBrowse := &cobra.Command{
		Use:   "browse [flags] year",
		Short: "",
		Long:  `TODO`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				// TODO check the cache location for what years we have
				// if we have no years, suggest running the
				// cache command
				return errors.New("requires a year (1950-2022)")
			}
			if _, err := strconv.Atoi(args[0]); err != nil {
				return fmt.Errorf("%q looks like it's not a year. Try 2021", args[0])
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			year, _ := strconv.Atoi(args[0])
			config := config.Init(cmd)
			bubble.Run(&config, int(year))
		},
	}

	return cmdBrowse
}
