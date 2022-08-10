package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type RawResponse [][][]string

func newRawResponse(r io.Reader) (RawResponse, error) {
	var res RawResponse
	return res, json.NewDecoder(r).Decode(&res)
}

func loadFromFile(filename string) (RawResponse, error) {
	jsonfile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonfile.Close()

	response, err := newRawResponse(jsonfile)
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

	response, err = newRawResponse(resp.Body)
	if err != nil {
		log.Printf("URL: %#v", url)
		log.Printf("Response: \n %#v", resp.Body)
	}

	if err := resp.Body.Close(); err != nil {
		log.Printf("Error closing http body: %v", err)
	}

	return response, err
}

func getAndParse(url UrlInfo) ([]RocketData, error) {
	response, err := get(url.Url)
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
