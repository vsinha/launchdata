package lib

import (
	"time"

	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
)

type PayloadData struct {
	Payload  string
	Operator string
	Orbit    string
	Function string
	Decay    string
	Outcome  string
	Cubesat  bool
}

type RocketData struct {
	TimestampRaw          string
	Timestamp             time.Time
	TBD                   bool
	Rocket                string
	FlightNumber          string
	LaunchSite            string
	LaunchServiceProvider string
	Notes                 string
	Payload               []PayloadData
}

type AllLaunchData struct {
	OrbitalFlights    []RocketData
	SuborbitalFlights []RocketData
}

var dryrun bool

// TODO Add a command for loading from file instead of doing an http request
func Command() {
	litter.Config.HomePackage = "lib"
	litter.Config.HidePrivateFields = false

	cmdCache := cmdCache()

	rootCmd := &cobra.Command{
		Use:   "launchdata",
		Short: "Launchdata 🚀\nA tool to download and examine rocket launch data from Wikipedia",
	}
	rootCmd.PersistentFlags().BoolVar(&dryrun, "dry-run", false, "Don't actually take any actions")

	rootCmd.AddCommand(cmdCache)
	rootCmd.Execute()
}
