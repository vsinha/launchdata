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
