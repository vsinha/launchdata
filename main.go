package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/sanity-io/litter"
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

type AllLauchData struct {
	orbitalFlights    []RocketData
	suborbitalFlights []RocketData
}

func shouldSkipEntry(entry []string) bool {
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

	return false
}

func parseTimestamp(raw string, year int) time.Time {
	// trim off the wiki link annotation
	raw = strings.Split(raw, "[")[0]
	// add the year
	raw = fmt.Sprintf("%d %s MST", year, raw)
	// there's no space between the month and the hour
	t, err := time.Parse("2006 2 January15:04:05 MST", raw)
	if err != nil {
		fmt.Printf("Error parsing time: %v\n", err)
	}
	if t.String() == "0001-01-01 00:00:00 +0000 UTC" {
		// try again without seconds, chinese launches seem to not
		// have seconds
		t, err = time.Parse("2006 2 January15:04 MST", raw)
		if err != nil {
			fmt.Printf("Error parsing time: %v\n", err)
		}
	}
	return t
}

func parseSingleDate(index *int, data [][]string, year int) (RocketData, error) {
	var rocketData RocketData
	var payloadData []PayloadData
	var i int

	// grab the date of the first entry
	date := data[*index][0]

	// Keep checking until the date changes
	for i = *index; i < len(data) && data[i][0] == date; i += 1 {
		if shouldSkipEntry(data[i]) {
			continue
		}

		if data[i][1] != "" {
			// The 2nd field (index 1) is always "" for payloads and notes, because they're indented
			// in the wiki table

			rocketData = RocketData{
				Datetime:     parseTimestamp(date, year),
				Rocket:       data[*index][1],
				FlightNumber: data[*index][3],
				LaunchSite:   data[*index][4],
				Lsp:          data[*index][6],
				Payload:      []PayloadData{},
			}
		} else if data[i][2] == data[i][3] && data[i][2] == data[i][5] {
			// The notes entry has the same piece of data represented in indexes 2-7,
			// so we can just check a couple of them
			rocketData.Notes = data[i][3]
		} else {
			data := PayloadData{
				Payload:  data[i][2],
				Operator: data[i][3],
				Orbit:    data[i][4],
				Function: data[i][5],
				Decay:    data[i][6],
				Outcome:  data[i][7],
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

	// The first 4 entries contain the month and some other header rows
	for index := 4; index < len(data); index++ {
		// fmt.Printf("Parsing:\n%s\n", litter.Sdump(data[index]))
		if shouldSkipEntry(data[index]) {
			continue
		}
		rocketData, err := parseSingleDate(&index, data, year)
		if err != nil {
			fmt.Println(err)
		}
		if rocketData.Payload == nil {
			fmt.Println(
				fmt.Errorf("Parsed a rocketData with nil payload, probably something has gone wrong:\n %s\n",
					litter.Sdump(rocketData)))
		}
		allRocketData = append(allRocketData, rocketData)
	}
	return allRocketData, nil
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

const baseUrl = "https://www.wikitable2json.com/api/List_of_spaceflight_launches_in"

type UrlInfo struct {
	Year int
	Url  string
}

// These wiki pages actually don't exist pre-2021
func generateUrlsForYearRange(startYear int, endYear int) []UrlInfo {
	var urls []UrlInfo

	for y := startYear; y <= endYear; y++ {
		urls = append(urls, UrlInfo{Year: y, Url: fmt.Sprintf("%s_January%%E2%%80%%93June_%d", baseUrl, y)})
		urls = append(urls, UrlInfo{Year: y, Url: fmt.Sprintf("%s_July%%E2%%80%%93December_%d", baseUrl, y)})
	}

	return urls
}

func getAndParse(url UrlInfo) ([]RocketData, error) {
	response, err := get(url.Url)
	// response, err := loadFromFile("launches-2022-jan-jun.json")
	if err != nil {
		return []RocketData{}, err
	}

	launchData, err := parseMultipleDates(response[0], url.Year)

	return launchData, err
}

func main() {
	urls := generateUrlsForYearRange(2021, 2022)
	for _, url := range urls {
		launchData, err := getAndParse(url)
		if err != nil {
			log.Fatalln(err)
			return
		}
		litter.Dump(launchData)

	}
}
