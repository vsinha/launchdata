package parse

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"launchdata/jsonio"
)

func TestCanLoadFile(t *testing.T) {
	response, err := jsonio.LoadFromFile("testdata/launches-2022-jan-jun.json")
	if err != nil {
		t.Errorf("%v", err)
	}
	if response == nil {
		t.Errorf("%v", response)
	}
}

func TestGeneratingUrlsForYearRange(t *testing.T) {
	got := generateUrlsForYearRange(2020, 2022)
	want := []UrlInfo{
		{
			Year:    2020,
			Url:     "https://www.wikitable2json.com/api/2020_in_spaceflight",
			WikiUrl: "https://en.wikipedia.org/wiki/2020_in_spaceflight",
		},
		{
			Year:    2021,
			Url:     "https://www.wikitable2json.com/api/List_of_spaceflight_launches_in_January%E2%80%93June_2021",
			WikiUrl: "https://en.wikipedia.org/wiki/List_of_spaceflight_launches_in_January%E2%80%93June_2021",
		},
		{
			Year:    2021,
			Url:     "https://www.wikitable2json.com/api/List_of_spaceflight_launches_in_July%E2%80%93December_2021",
			WikiUrl: "https://en.wikipedia.org/wiki/List_of_spaceflight_launches_in_July%E2%80%93December_2021",
		},
		{
			Year:    2022,
			Url:     "https://www.wikitable2json.com/api/List_of_spaceflight_launches_in_January%E2%80%93June_2022",
			WikiUrl: "https://en.wikipedia.org/wiki/List_of_spaceflight_launches_in_January%E2%80%93June_2022",
		},
		{
			Year:    2022,
			Url:     "https://www.wikitable2json.com/api/List_of_spaceflight_launches_in_July%E2%80%93December_2022",
			WikiUrl: "https://en.wikipedia.org/wiki/List_of_spaceflight_launches_in_July%E2%80%93December_2022",
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("diff (+want,-got:\n%s", diff)
	}
}

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
