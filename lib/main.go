package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
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

var dryrun bool

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
	var response RawResponse

	if dryrun {
		fmt.Printf("Dry run: Would request %s\n", url)
		return response, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

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

		fmt.Printf("Parsed %d orbital launches in %d (%s, %s)\n", len(launchData), url.Year, url.Url, url.WikiUrl)
		allLaunchData.OrbitalFlights = append(allLaunchData.OrbitalFlights, launchData...)
	}
	return allLaunchData, nil
}

func toJson(contents interface{}) (*bytes.Buffer, error) {
	res, err := json.Marshal(contents)
	if err != nil {
		return nil, err
	}

	formattedJson := &bytes.Buffer{}
	if err := json.Indent(formattedJson, []byte(res), "", "  "); err != nil {
		return nil, err
	}
	return formattedJson, nil
}

func writeJsonFile(contents interface{}, filename string) error {
	formattedJson, err := toJson(contents)
	if err != nil {
		return err
	}

	if dryrun {
		fmt.Printf("Dry run: would output file %s\n", filename)
		return nil
	}

	if err := os.WriteFile(filename, formattedJson.Bytes(), 0o644); err != nil {
		return err
	}

	return nil
}

func getAndWrite(startYear int, endYear int, filename string) {
	if dryrun {
		fmt.Printf("Dry run: would get and write file %s\n", filename)
		return
	}

	results, _ := getAndParseMultipleYears(startYear, endYear)

	if filename != "" {
		fmt.Printf("Writing %s\n", filename)
		if err := writeJsonFile(results, filename); err != nil {
			panic(err)
		}
	}
}

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
				getAndWrite(i, i, filename)
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
				getAndWrite(startYear, endYear, outputFilename)
			} else {
				getAndWrite(year, year, outputFilename)
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
