package parse

import (
	"fmt"
	"strings"
	"time"
)

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
	err := fmt.Errorf("no parser user") // <-- error strings should be lower case https://github.com/golang/go/wiki/CodeReviewComments#error-strings
	for _, format := range formats {
		t, err = time.Parse(format, raw)
		if err == nil {
			break
		}
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
