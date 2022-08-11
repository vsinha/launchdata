package parse

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"launchdata/config"
	"launchdata/jsonio"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/sanity-io/litter"
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

func LoadLaunchDataFromFile(filename string) (AllLaunchData, error) {
	jsonfile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonfile.Close()

	var response AllLaunchData
	json.NewDecoder(jsonfile).Decode(&response)
	return response, err
}

var months mapset.Set[string] = mapset.NewSet("January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December")

func shouldSkipEntry(entry []string) bool {
	if len(entry) == 0 {
		return true
	}
	if len(entry) == 2 {
		// this entry only contains the date and ""
		return true
	}
	if r, _ := utf8.DecodeRune([]byte(entry[0])); r == '←' {
		// "←  Jan\nFeb\nMar\nApr\nMay\nJun\nJul\nAug\nSep\nOct\nNov\nDec →"
		return true
	}
	if strings.HasPrefix(entry[0], "For flights after") {
		return true
	}
	if months.Contains(entry[0]) && months.Contains(entry[1]) {
		return true
	}
	if strings.Contains(entry[3], "Upcoming launches") {
		return true
	}

	return false
}

var wikiLinkRemovalRegex = regexp.MustCompile(`\[\d+\]`)

func cleanWikilink(input string) string {
	return wikiLinkRemovalRegex.ReplaceAllString(input, "")
}

func checkIfCubesat(input string) (string, bool) {
	isCubesat := false
	if strings.ContainsAny(input, "⚀▫") {
		isCubesat = true
		input = strings.ReplaceAll(input, "⚀", "")
		input = strings.ReplaceAll(input, "▫", "")
		input = strings.TrimSpace(input)
	}
	return input, isCubesat
}

func parseSingleDate(index *int, data [][]string, year int) (RocketData, error) {
	var rocketData RocketData
	var payloadData []PayloadData
	var i int

	if len(data[*index]) < 1 {
		return rocketData, fmt.Errorf("no data for year %d, index %d", year, index)
	}

	// grab the timestampRaw of the first entry
	timestampRaw := data[*index][0]

	// Keep checking until the date changes
	for i = *index; i < len(data) && len(data[i]) > 0 && data[i][0] == timestampRaw; i += 1 {

		if shouldSkipEntry(data[i]) {
			continue
		}

		timestampRawCleaned := cleanWikilink(timestampRaw)
		for j, entry := range data[i] {
			data[i][j] = cleanWikilink(entry)
		}

		if data[i][1] != "" {
			// The 2nd field (index 1) is always "" for payloads and notes, because they're indented
			// in the wiki table

			var timestamp time.Time
			var tbd bool
			if strings.Contains(timestampRaw, "TBD") {
				tbd = true
			} else {
				var err error
				timestamp, err = parseTimestamp(timestampRawCleaned, year)
				if err != nil {
					fmt.Println(fmt.Errorf("failed to parse timestamp %v", err))
				}
			}

			if len(data[i]) < 7 {
				continue
			}

			rocketData = RocketData{
				TimestampRaw:          timestampRaw,
				Timestamp:             timestamp,
				TBD:                   tbd,
				Rocket:                data[*index][1],
				FlightNumber:          data[*index][3],
				LaunchSite:            data[*index][4],
				LaunchServiceProvider: data[*index][6],
				Payload:               []PayloadData{},
			}
		} else if data[i][2] == data[i][3] && data[i][2] == data[i][5] {
			// The notes entry has the same piece of data represented in indexes 2-7,
			// so we can just check a couple of them
			rocketData.Notes = data[i][3]
		} else {

			if len(data[i]) < 8 {
				continue
			}

			payload, cubesat := checkIfCubesat(data[i][2])
			data := PayloadData{
				Payload:  payload,
				Operator: data[i][3],
				Orbit:    data[i][4],
				Function: data[i][5],
				Decay:    data[i][6],
				Outcome:  data[i][7],
				Cubesat:  cubesat,
			}
			payloadData = append(payloadData, data)
		}
	}
	*index = i - 1

	rocketData.Payload = payloadData
	return rocketData, nil
}

func parseMultipleDates(data [][]string, year int) ([]RocketData, error) {
	var allRocketData []RocketData
	now := time.Now()

	// The first 4 entries contain the month and some other header rows
	for index := 4; index < len(data); index++ {
		if shouldSkipEntry(data[index]) {
			continue
		}

		rocketData, err := parseSingleDate(&index, data, year)
		if err != nil {
			fmt.Println(err)
		}

		// Did it launch empty?
		if rocketData.Payload == nil && !rocketData.TBD && rocketData.Timestamp.Before(now) {
			fmt.Println(rocketData.Timestamp)
			fmt.Println(
				fmt.Errorf("parsed a rocketData with no payload, probably something has gone wrong:\n %s",
					litter.Sdump(rocketData)))
		}

		allRocketData = append(allRocketData, rocketData)
	}

	return allRocketData, nil
}

func getAndParse(config config.Config, url UrlInfo) ([]RocketData, error) {
	response, err := jsonio.Get(config, url.Url)
	if err != nil {
		return []RocketData{}, err
	}

	launchData, err := parseMultipleDates(response[0], url.Year)

	return launchData, err
}

func getAndParseMultipleYears(config config.Config, startYear int, endYear int) (AllLaunchData, error) {
	urls := generateUrlsForYearRange(startYear, endYear)
	var allLaunchData AllLaunchData
	for _, url := range urls {
		launchData, err := getAndParse(config, url)
		if err != nil {
			fmt.Printf("Encountered an error: %v\n", err)
		}

		fmt.Printf("Parsed %d orbital launches in %d (%s, %s)\n", len(launchData), url.Year, url.Url, url.WikiUrl)
		allLaunchData.OrbitalFlights = append(allLaunchData.OrbitalFlights, launchData...)
	}
	return allLaunchData, nil
}

func GetAndWrite(config config.Config, startYear int, endYear int, filename string) {
	if config.DryRun {
		fmt.Printf("Dry run: would get and write file %s\n", filename)
		return
	}

	results, _ := getAndParseMultipleYears(config, startYear, endYear)

	if filename != "" {
		fmt.Printf("Writing %s\n", filename)
		if err := jsonio.WriteJsonFile(config, results, filename); err != nil {
			panic(err)
		}
	}
}
