package jsonio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"launchdata/config"
)

type RawResponse [][][]string

func newRawResponse(r io.Reader) (RawResponse, error) {
	var res RawResponse
	return res, json.NewDecoder(r).Decode(&res)
}

func LoadFromFile(filename string) (RawResponse, error) {
	jsonfile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonfile.Close()

	response, err := newRawResponse(jsonfile)
	return response, err
}

func Get(config config.Config, url string) (RawResponse, error) {
	var response RawResponse

	if config.DryRun {
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

func WriteJsonFile(config config.Config, contents interface{}, filename string) error {
	formattedJson, err := toJson(contents)
	if err != nil {
		return err
	}

	if config.DryRun {
		fmt.Printf("Dry run: would output file %s\n", filename)
		return nil
	}

	if err := os.WriteFile(filename, formattedJson.Bytes(), 0o644); err != nil {
		return err
	}

	return nil
}
