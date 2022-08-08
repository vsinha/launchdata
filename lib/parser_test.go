package lib

import (
	"encoding/json"
	"testing"

	golden "github.com/jimeh/go-golden"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCleanWikilink(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1 March21:38:00[97]", "1 March21:38:00"},
		{"foo[0][2]", "foo"},
		{"foo[hello][2]", "foo[hello]"},
	}

	for _, test := range tests {
		got := cleanWikilink(test.input)
		if !(got == test.want) {
			t.Errorf("wanted: %v, got: %v", test.want, got)
		}
	}
}

func TestSkippingIrrelevantEntries(t *testing.T) {
	data := [][]string{
		{
			"←  Jan\nFeb\nMar\nApr\nMay\nJun\nJul\nAug\nSep\nOct\nNov\nDec →\n\n\nMarch",
			"←  Jan\nFeb\nMar\nApr\nMay\nJun\nJul\nAug\nSep\nOct\nNov\nDec →\n\n\nMarch",
			"←  Jan\nFeb\nMar\nApr\nMay\nJun\nJul\nAug\nSep\nOct\nNov\nDec →\n\n\nMarch",
			"←  Jan\nFeb\nMar\nApr\nMay\nJun\nJul\nAug\nSep\nOct\nNov\nDec →\n\n\nMarch",
			"←  Jan\nFeb\nMar\nApr\nMay\nJun\nJul\nAug\nSep\nOct\nNov\nDec →\n\n\nMarch",
			"←  Jan\nFeb\nMar\nApr\nMay\nJun\nJul\nAug\nSep\nOct\nNov\nDec →\n\n\nMarch",
			"←  Jan\nFeb\nMar\nApr\nMay\nJun\nJul\nAug\nSep\nOct\nNov\nDec →\n\n\nMarch",
			"←  Jan\nFeb\nMar\nApr\nMay\nJun\nJul\nAug\nSep\nOct\nNov\nDec →\n\n\nMarch",
			"←  Jan\nFeb\nMar\nApr\nMay\nJun\nJul\nAug\nSep\nOct\nNov\nDec →",
		},
		{"←  Jan\nFeb\nMar\nApr\nMay\nJun\nJul\nAug\nSep\nOct\nNov\nDec →"},
		{
			"1 March21:38:00[97]",
			"Atlas V 541",
			"Atlas V 541",
			"AV-095",
			"Cape Canaveral SLC-41",
			"Cape Canaveral SLC-41",
			"ULA",
			"ULA",
		},
		{
			"1 March21:38:00[97]",
			"",
		},
	}

	var filtered [][]string
	for _, entry := range data {
		if !shouldSkipEntry(entry) {
			filtered = append(filtered, entry)
		}
	}

	if len(filtered) != 1 {
		t.Errorf("Failed to filter irrelevant entries: %v\n", filtered)
	}
}

func TestCanParseSingleDateWithSinglePayload(t *testing.T) {
	response, err := loadFromFile("testdata/launches-2022-jan-6.json")
	require.NoError(t, err)

	index := 0
	got, err := parseSingleDate(&index, response[0], 2022)
	require.NoError(t, err)

	gotJson, err := json.Marshal(&got)
	require.NoError(t, err)

	if golden.Update() {
		golden.Set(t, gotJson)
	}
	want := golden.Get(t)

	assert.Equal(t, want, gotJson)
}

func TestCanParseSingleDateWithMultiplePayloads(t *testing.T) {
	response, err := loadFromFile("testdata/launches-2022-jan-13.json")
	require.NoError(t, err)
	index := 0

	got, err := parseSingleDate(&index, response[0], 2022)
	require.NoError(t, err)

	verify(t, got)
}

func TestCanParseMultipleDates(t *testing.T) {
	response, err := loadFromFile("testdata/launches-2022-jan-6-17.json")
	require.NoError(t, err)

	got, err := parseMultipleDates(response[0], 2022)
	require.NoError(t, err)

	verify(t, got)
}

func TestCanParseFullWikiPage(t *testing.T) {
	response, err := loadFromFile("testdata/launches-2022-jan-jun.json")
	require.NoError(t, err)

	got, err := parseMultipleDates(response[0], 2022)
	require.NoError(t, err)

	verify(t, got)
}
