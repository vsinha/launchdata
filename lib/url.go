package lib

import (
	"fmt"
)

const (
	baseUrl     = "https://www.wikitable2json.com/api"
	baseWikiUrl = "https://en.wikipedia.org/wiki"
)

type UrlInfo struct {
	Year    int
	Url     string
	WikiUrl string
}

// These wiki pages actually don't exist pre-2021
func generateUrlsForYearRange(startYear int, endYear int) []UrlInfo {
	var urls []UrlInfo

	for y := startYear; y <= endYear; y++ {
		if y < 2021 {
			urls = append(urls, UrlInfo{
				Year:    y,
				Url:     fmt.Sprintf("%s/%d_in_spaceflight", baseUrl, y),
				WikiUrl: fmt.Sprintf("%s/%d_in_spaceflight", baseWikiUrl, y),
			})
		} else {
			urls = append(urls, UrlInfo{
				Year:    y,
				Url:     fmt.Sprintf("%s/List_of_spaceflight_launches_in_January%%E2%%80%%93June_%d", baseUrl, y),
				WikiUrl: fmt.Sprintf("%s/List_of_spaceflight_launches_in_January%%E2%%80%%93June_%d", baseWikiUrl, y),
			})
			urls = append(urls, UrlInfo{
				Year:    y,
				Url:     fmt.Sprintf("%s/List_of_spaceflight_launches_in_July%%E2%%80%%93December_%d", baseUrl, y),
				WikiUrl: fmt.Sprintf("%s/List_of_spaceflight_launches_in_July%%E2%%80%%93December_%d", baseWikiUrl, y),
			})
		}
	}

	return urls
}
