package parse

import (
	"testing"
)

func TestParsingDatestamp(t *testing.T) {
	tests := []struct {
		input string
		year  int
		want  string
	}{
		{"13 January15:25:39[2]", 2021, "2021-01-13 15:25:39 +0000 MST"},
		{"13 January22:51:39[46][47]", 2021, "2021-01-13 22:51:39 +0000 MST"},
		{"21 January19:00:00[53]", 2022, "2022-01-21 19:00:00 +0000 MST"},
		{"2022 8 March~05:06 MST", 2022, "2022-03-08 05:06:00 +0000 MST"},
	}

	for _, test := range tests {
		got, _ := parseTimestamp(test.input, test.year)
		want := timeParse(test.want)
		if !got.Equal(want) {
			t.Errorf("wanted: %v, got: %v", want, got)
		}
	}
}
