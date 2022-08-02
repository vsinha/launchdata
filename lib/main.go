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

func getAndParseMultipleYears(startYear int, endYear int) (AllLaunchData, error) {
	urls := generateUrlsForYearRange(startYear, endYear)
	var allLaunchData AllLaunchData
	for _, url := range urls {
		launchData, err := getAndParse(url)
		if err != nil {
			fmt.Printf("Encountered an error: %v\n", err)
		}

		fmt.Printf("Parsed %d orbital launches in %d (%s)\n", len(launchData), url.Year, url.Url)
		allLaunchData.OrbitalFlights = append(allLaunchData.OrbitalFlights, launchData...)
	}
	return allLaunchData, nil
}

func writeJsonFile(contents interface{}, filename string) error {
	res, err := json.Marshal(contents)
	if err != nil {
		return err
	}

	formattedJson := &bytes.Buffer{}
	if err := json.Indent(formattedJson, []byte(res), "", "  "); err != nil {
		return err
	}

	if err := os.WriteFile(filename, formattedJson.Bytes(), 0o644); err != nil {
		return err
	}

	return nil
}

// TODO Add a command for loading from file instead of doing an http request
func Command() {
	litter.Config.HomePackage = "lib"
	litter.Config.HidePrivateFields = false

	var startYear int
	var endYear int
	var outputFilename string

	rootCmd := &cobra.Command{
		Use:   "launchdata",
		Short: "Launchdata - a simple CLI to transform and inspect strings",
		Long:  `TODO`,
		Run: func(cmd *cobra.Command, args []string) {
			results, _ := getAndParseMultipleYears(startYear, endYear)

			if outputFilename != "" {
				if err := writeJsonFile(results, outputFilename); err != nil {
					panic(err)
				}
			}
		},
	}
	rootCmd.Flags().IntVarP(&startYear, "start", "s", 2021, "Start Year")
	rootCmd.MarkFlagRequired("start")

	rootCmd.Flags().IntVarP(&endYear, "end", "e", 2021, "End Year")
	rootCmd.MarkFlagRequired("end")

	rootCmd.Flags().StringVarP(&outputFilename, "output", "o", "", "JSON output file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
