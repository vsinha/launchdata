package parse

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type TimeData struct {
	TimestampRaw   string
	TimestampClean string
	Timestamp      time.Time
	Tbd            bool
	ParsedOk       bool
	ParseErr       error
}

func (t TimeData) LaunchedAlready(now time.Time) bool {
	return !t.Tbd && t.ParsedOk && t.Timestamp.Before(now)
}

func (t TimeData) DateString() string {
	if !t.ParsedOk {
		return t.TimestampClean
	}
	return t.Timestamp.Format("2006-01-02")
}

func (t TimeData) TimeString() string {
	if !t.ParsedOk {
		return t.TimestampClean
	}
	return t.Timestamp.Format("2006-01-02 15:04 (UTC)")
}

func parseTimestampFormat(raw string, year int) (time.Time, error) {
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
	err := fmt.Errorf("no parser user") // <-- error strings should be lower case https://github.com/golang/go/wiki/CodeReviewComments#error-strings
	for _, format := range formats {
		t, err = time.Parse(format, raw)
		if err == nil {
			break
		}
	}

	// Future timestamp strings will contain these terms, we can simply return no error, as it's OK to have failed to parse them
	if err != nil &&
		(strings.Contains(raw, "Early") || strings.Contains(raw, "Mid") || strings.Contains(raw, "Late")) {
		return t, errors.New("failed to parse, contains early/mid/late")
	}

	if t.String() == "0001-01-01T00:00:00Z" {
		return t, errors.New("failed to parse, resulted in default date")
	}

	return t, err
}

func parseTimestamp(raw string, year int) TimeData {
	var timestamp time.Time
	var tbd bool
	var err error

	cleaned := cleanWikilink(raw)

	if strings.Contains(raw, "TBD") {
		tbd = true
		err = errors.New("TBD")
	} else {
		timestamp, err = parseTimestampFormat(cleaned, year)
	}

	time := TimeData{
		TimestampRaw:   raw,
		TimestampClean: cleaned,
		Timestamp:      timestamp,
		Tbd:            tbd,
		ParsedOk:       err == nil,
		ParseErr:       err,
	}

	return time
}
