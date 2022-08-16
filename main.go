package main

import (
	"launchdata/cmd"
)

func main() {
	rootCmd := cmd.Root()
	rootCmd.Execute()
	// if err := rootCmd.Execute(); err != nil {
	// 	panic(err)
	// }
}
