package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	// "github.com/sanity-io/litter"
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
}

type RocketData struct {
	Datetime     time.Time
	Rocket       string
	FlightNumber string
	LaunchSite   string
	Lsp          string
	Notes        string
	Payload      []PayloadData
}

type AllLaunchData struct {
	orbitalFlights    []RocketData
	suborbitalFlights []RocketData
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
	urls := generateUrlsForYearRange(2022, 2022)
	var allLaunchData AllLaunchData
	for _, url := range urls {
		launchData, err := getAndParse(url)
		if err != nil {
			return allLaunchData, err
		}
		// litter.Dump(launchData)
		fmt.Printf("Parsed %d orbital launches in %d", len(launchData), url.Year)
		allLaunchData.orbitalFlights = append(allLaunchData.orbitalFlights, launchData...)
	}
	return allLaunchData, nil
}

func Command() {
	var startYear int
	var endYear int

	rootCmd := &cobra.Command{
		Use:   "launchdata",
		Short: "Launchdata - a simple CLI to transform and inspect strings",
		Long:  `TODO`,
		Run: func(cmd *cobra.Command, args []string) {
			getAndParseMultipleYears(startYear, endYear)
		},
	}
	rootCmd.Flags().IntVarP(&startYear, "start", "s", 2022, "Start Year")
	rootCmd.MarkFlagRequired("start")

	rootCmd.Flags().IntVarP(&endYear, "end", "e", 2022, "End Year")
	rootCmd.MarkFlagRequired("end")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
