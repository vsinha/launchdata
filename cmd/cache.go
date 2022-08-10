package cmd

import (
	"fmt"
	"path"

	"launchdata/parse"

	"github.com/spf13/cobra"
)

func cmdCacheAll() *cobra.Command {
	var outputDir string
	cmdCacheAll := &cobra.Command{
		Use:   "all",
		Short: "Download all historical launch data from wikipedia",
		Long:  `TODO`,
		Run: func(cmd *cobra.Command, args []string) {
			from := 1951
			to := 2022
			fmt.Printf("Caching all files from %d to %d\n", from, to)
			for i := from; i <= to; i++ {
				filename := path.Join(outputDir, fmt.Sprintf("launchdata-%d.json", i))
				parse.GetAndWrite(i, i, filename)
			}
		},
	}
	cmdCacheAll.Flags().StringVar(&outputDir, "output-dir", "./data", "output directory")

	return cmdCacheAll
}

func cmdCache() *cobra.Command {
	var year int
	var startYear int
	var endYear int
	var outputFilename string

	cmdCache := &cobra.Command{
		Use:   "cache",
		Short: "Download launch data from wikipedia and cache it locally",
		Long:  `TODO`,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flags().Changed("startYear") {
				parse.GetAndWrite(startYear, endYear, outputFilename)
			} else {
				parse.GetAndWrite(year, year, outputFilename)
			}
		},
	}
	cmdCache.Flags().IntVarP(&startYear, "start", "s", 2021, "Start Year")
	cmdCache.Flags().IntVarP(&endYear, "end", "e", 2021, "End Year")
	cmdCache.MarkFlagsRequiredTogether("start", "end")

	cmdCache.Flags().IntVarP(&year, "year", "y", 2021, "Specify a single year")
	cmdCache.MarkFlagsMutuallyExclusive("year", "start")

	cmdCache.Flags().StringVarP(&outputFilename, "output", "o", "", "JSON output file")

	cmdCacheAll := cmdCacheAll()
	cmdCache.AddCommand(cmdCacheAll)

	return cmdCache
}
