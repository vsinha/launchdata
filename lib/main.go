package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
)

type Response struct {
	UserId int    `json:"userId"`
	Body   string `json:"body"`
	Title  string `json:"title"`
}

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

type RawResponse [][][]string

func (res *RawResponse) decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&res)
}

func get(url string) (RawResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var response RawResponse
	err = response.decode(resp.Body)
	if err != nil {
		log.Printf("URL: %#v", url)
		log.Printf("Response: \n %#v", resp.Body)
	}
	return response, err
}

func loadFromFile(filename string) (RawResponse, error) {
	jsonfile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonfile.Close()

	var response RawResponse
	err = response.decode(jsonfile)
	return response, err
}

func getAndParse(url UrlInfo) ([]RocketData, error) {
	response, err := get(url.Url)
	// response, err := loadFromFile("launches-2022-jan-jun.json")
	// litter.Dump(response)
	if err != nil {
		return []RocketData{}, err
	}

	launchData, err := parseMultipleDates(response[0], url.Year)

	return launchData, err
}

// TODO add support for outputting JSON

func getAndParseMultipleYears(startYear int, endYear int) (AllLaunchData, error) {
	urls := generateUrlsForYearRange(startYear, endYear)
	var allLaunchData AllLaunchData
	for _, url := range urls {
		launchData, err := getAndParse(url)
		if err != nil {
			fmt.Printf("Encountered an error: %v\n", err)
		}
		// litter.Dump(launchData)
		fmt.Printf("Parsed %d orbital launches in %d (%s)\n", len(launchData), url.Year, url.Url)
		allLaunchData.OrbitalFlights = append(allLaunchData.OrbitalFlights, launchData...)
	}
	// fmt.Printf("%#v\n", allLaunchData)
	return allLaunchData, nil
}

func Command() {
	litter.Config.HomePackage = "lib"
	litter.Config.HidePrivateFields = false

	var startYear int
	var endYear int
	var output string

	rootCmd := &cobra.Command{
		Use:   "launchdata",
		Short: "Launchdata - a simple CLI to transform and inspect strings",
		Long:  `TODO`,
		Run: func(cmd *cobra.Command, args []string) {
			results, _ := getAndParseMultipleYears(startYear, endYear)

			if output != "" {
				resultsJson, err := json.Marshal(results)
				if err != nil {
					panic(err)
				}

				formattedJson := &bytes.Buffer{}
				if err := json.Indent(formattedJson, []byte(resultsJson), "", "  "); err != nil {
					panic(err)
				}

				if err := os.WriteFile(output, formattedJson.Bytes(), 0o644); err != nil {
					panic(err)
				}
			}
		},
	}
	rootCmd.Flags().IntVarP(&startYear, "start", "s", 2021, "Start Year")
	rootCmd.MarkFlagRequired("start")

	rootCmd.Flags().IntVarP(&endYear, "end", "e", 2021, "End Year")
	rootCmd.MarkFlagRequired("end")

	rootCmd.Flags().StringVarP(&output, "output", "o", "", "JSON output file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
