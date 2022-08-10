package config

import "github.com/spf13/cobra"

type Config struct {
	DryRun bool
}

func Init(cmd *cobra.Command) Config {
	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		panic(err)
	}

	return Config{
		DryRun: dryRun,
	}
}
