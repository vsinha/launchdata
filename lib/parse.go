package lib

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/sanity-io/litter"

	mapset "github.com/deckarep/golang-set/v2"
)

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

// TODO wrap the Time in a struct which includes an 'ok' bool,
// to figure out whether we parsed it correctly or not.
// I'd use Option[Time] but this is Go
func parseTimestamp(raw string, year int) (time.Time, error) {
	raw = cleanWikilink(raw)
	raw = strings.TrimSpace(raw)

	// add the year
	if !strings.HasPrefix(raw, fmt.Sprint(year)) {
		raw = fmt.Sprintf("%d %s", year, raw)
	}

	formats := []string{
		"2006 2 January~15:04 MST",
		"2006 2 January15:04 (UTC)",
		"2006 2 January15:04:05 (UTC)",
		"2006 2 January~15:04",
		"2006 2 January15:04:05",
		"2006 2 January15:04",
		"2006 2 January",
	}

	var t time.Time
	err := fmt.Errorf("Haven't attempted any parsing yet")
	for i := 0; i < len(formats) && err != nil; i++ {
		t, err = time.Parse(formats[i], raw)
	}

	// TODO add a check to see if the time is still "0001-01-01T00:00:00Z", log
	// an error or something if it is

	// Future timestamp strings will contain these terms, we can simply return
	// no error, as it's OK to have failed to parse them
	if err != nil &&
		(strings.Contains(raw, "Early") || strings.Contains(raw, "Mid") || strings.Contains(raw, "Late")) {
		return t, nil
	}

	return t, err
}

func cleanWikilink(input string) string {
	m := regexp.MustCompile(`\[\d+\]`)
	return m.ReplaceAllString(input, "")
}

func checkIfCubesat(input string) (string, bool) {
	isCubesat := false
	if strings.ContainsAny(input, "⚀▫") {
		isCubesat = true
		input = strings.Replace(input, "⚀", "", 1)
		input = strings.Replace(input, "▫", "", 1)
		input = strings.TrimSpace(input)
	}
	return input, isCubesat
}

func parseSingleDate(index *int, data [][]string, year int) (RocketData, error) {
	var rocketData RocketData
	var payloadData []PayloadData
	var i int

	// grab the timestampRaw of the first entry
	timestampRaw := data[*index][0]
	// date = cleanWikilink(date)
	// fmt.Printf("%#v\n", date)

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
					fmt.Println(fmt.Errorf("Failed to parse timestamp %v", err))
				}
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
		// fmt.Printf("Parsing:\n%s\n", litter.Sdump(data[index]))
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
				fmt.Errorf("Parsed a rocketData with no payload, probably something has gone wrong:\n %s\n",
					litter.Sdump(rocketData)))
		}

		allRocketData = append(allRocketData, rocketData)
	}

	return allRocketData, nil
}
