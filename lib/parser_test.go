package lib

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCleanWikilink(t *testing.T) {
	tests := []struct {
		test string
		want string
	}{
		{test: "1 March21:38:00[97]", want: "1 March21:38:00"},
		{test: "foo[0][2]", want: "foo"},
		{test: "foo[hello][2]", want: "foo[hello]"},
	}

	for _, test := range tests {
		got := cleanWikilink(test.test)
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
	index := 0
	got, err := parseSingleDate(&index, response[0], 2022)
	if err != nil {
		t.Errorf("parse error: %v", err)
	}

	want := RocketData{
		// Datetime:     "6 January21:49:10[1]",
		Datetime:              timeParse("2022-01-06 21:49:10 +0000 MST"),
		Rocket:                "Falcon 9 Block 5",
		FlightNumber:          "Starlink Group 4-5",
		LaunchSite:            "Kennedy LC-39A",
		LaunchServiceProvider: "SpaceX",
		Payload: []PayloadData{
			{
				Payload:  "Starlink × 49",
				Operator: "SpaceX",
				Orbit:    "Low Earth",
				Function: "Communications",
				Decay:    "In orbit",
				Outcome:  "Operational",
			},
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("diff (-got,+want:\n%s", diff)
	}
}

func TestCanParseSingleDateWithMultiplePayloads(t *testing.T) {
	response, err := loadFromFile("testdata/launches-2022-jan-13.json")
	index := 0
	got, err := parseSingleDate(&index, response[0], 2022)
	if err != nil {
		t.Errorf("parse error: %v", err)
	}

	want := RocketData{
		Datetime:              timeParse("2022-01-13 15:25:39 +0000 MST"),
		Rocket:                "Falcon 9 Block 5",
		FlightNumber:          "Transporter-3",
		LaunchSite:            "Cape Canaveral SLC-40",
		LaunchServiceProvider: "SpaceX",
		Notes:                 "Dedicated SmallSat Rideshare mission to sun-synchronous orbit, designated Transporter-3.",
		Payload: []PayloadData{
			{
				Payload:  "ION SCV-004 Elysian Eleonora",
				Operator: "D-Orbit",
				Orbit:    "Low Earth (SSO)",
				Function: "CubeSat deployer",
				Decay:    "In orbit",
				Outcome:  "Operational",
			},
			{
				Payload:  "Alba Cluster 3That time of year",
				Operator: "Alba Orbital",
				Orbit:    "Low Earth (SSO)",
				Function: "PocketQube dispenser",
				Decay:    "In orbit",
				Outcome:  "Operational",
			},
			{
				Payload:  "Alba Cluster 4",
				Operator: "Alba Orbital",
				Orbit:    "Low Earth (SSO)",
				Function: "PocketQube dispenser",
				Decay:    "In orbit",
				Outcome:  "Operational",
			},
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("diff (-got,+want:\n%s", diff)
	}
}

func TestCanParseMultipleDates(t *testing.T) {
	response, err := loadFromFile("testdata/launches-2022-jan-6-17.json")
	got, err := parseMultipleDates(response[0], 2022)
	if err != nil {
		t.Errorf("parse error: %v", err)
	}

	want := []RocketData{
		{
			Datetime:              timeParse("2022-01-06 21:49:10 +0000 MST"),
			Rocket:                "Falcon 9 Block 5",
			FlightNumber:          "Starlink Group 4-5",
			LaunchSite:            "Kennedy LC-39A",
			LaunchServiceProvider: "SpaceX",
			Notes:                 "",
			Payload: []PayloadData{
				{
					Payload:  "Starlink × 49",
					Operator: "SpaceX",
					Orbit:    "Low Earth",
					Function: "Communications",
					Decay:    "In orbit",
					Outcome:  "Operational",
				},
			},
		},
		{
			Datetime:              timeParse("2022-01-13 15:25:39 +0000 MST"),
			Rocket:                "Falcon 9 Block 5",
			FlightNumber:          "Transporter-3",
			LaunchSite:            "Cape Canaveral SLC-40",
			LaunchServiceProvider: "SpaceX",
			Notes:                 "Dedicated SmallSat Rideshare mission to sun-synchronous orbit, designated Transporter-3.",
			Payload: []PayloadData{
				{
					Payload:  "ION SCV-004 Elysian Eleonora",
					Operator: "D-Orbit",
					Orbit:    "Low Earth (SSO)",
					Function: "CubeSat deployer",
					Decay:    "In orbit",
					Outcome:  "Operational",
				},
				{
					Payload:  "Alba Cluster 3That time of year",
					Operator: "Alba Orbital",
					Orbit:    "Low Earth (SSO)",
					Function: "PocketQube dispenser",
					Decay:    "In orbit",
					Outcome:  "Operational",
				},
			},
		},
		{
			Datetime:              timeParse("2022-01-13 22:51:39 +0000 MST"),
			Rocket:                "LauncherOne",
			FlightNumber:          "\"Above the Clouds\"",
			LaunchSite:            "Cosmic Girl, Mojave",
			LaunchServiceProvider: "Virgin Orbit",
			Notes:                 "STP-27VPB mission (ELaNa 29, GEARRS-3, and TechEdSat-3) for the Defense Innovation Unit. The ELaNa 29 mission consists of two CubeSats (PAN-A and PAN-B) that will autonomously rendezvous and dock in low Earth orbit.",
			Payload: []PayloadData{
				{
					Payload:  "Lemur-2-Krywe (ADLER-1)",
					Operator: "Austrian Space Forum",
					Orbit:    "Low Earth",
					Function: "Space debris measurement",
					Decay:    "In orbit",
					Outcome:  "Operational",
					Cubesat:  true,
				},
				{
					Payload:  "GEARRS-3",
					Operator: "Air Force Research Center",
					Orbit:    "Low Earth",
					Function: "Technology demonstration",
					Decay:    "In orbit",
					Outcome:  "Operational",
					Cubesat:  true,
				},
			},
		},

		{
			Datetime:              timeParse("2022-01-17 02:35:00 +0000 MST"),
			Rocket:                "Long March 2D",
			FlightNumber:          "2D-Y70",
			LaunchSite:            "Taiyuan LC-9",
			LaunchServiceProvider: "CASC",
			Notes:                 "",
			Payload: []PayloadData{
				{
					Payload:  "Shiyan-13",
					Operator: "CAS",
					Orbit:    "Low Earth (SSO)",
					Function: "Technology demonstration",
					Decay:    "In orbit",
					Outcome:  "Operational",
				},
			},
		},
		{
			Datetime:              timeParse("2022-01-19 02:02:40 +0000 MST"),
			Rocket:                "Falcon 9 Block 5",
			FlightNumber:          "Starlink Group 4-6",
			LaunchSite:            "Kennedy LC-39A",
			LaunchServiceProvider: "SpaceX",
			Notes:                 "",
			Payload: []PayloadData{
				{
					Payload:  "Starlink × 49",
					Operator: "SpaceX",
					Orbit:    "Low Earth",
					Function: "Communications",
					Decay:    "In orbit",
					Outcome:  "Operational",
				},
			},
		},
		{
			Datetime:              timeParse("2022-01-21 19:00:00 +0000 MST"),
			Rocket:                "Atlas V 511",
			FlightNumber:          "AV-084",
			LaunchSite:            "Cape Canaveral SLC-41",
			LaunchServiceProvider: "ULA",
			Notes:                 "First and only flight of the 511 configuration for Atlas V.",
			Payload: []PayloadData{
				{
					Payload:  "USSF-8 / GSSAP-5",
					Operator: "U.S. Space Force",
					Orbit:    "Geosynchronous",
					Function: "Space surveillance",
					Decay:    "In orbit",
					Outcome:  "Operational",
				},
				{
					Payload:  "USSF-8 / GSSAP-6",
					Operator: "U.S. Space Force",
					Orbit:    "Geosynchronous",
					Function: "Space surveillance",
					Decay:    "In orbit",
					Outcome:  "Operational",
				},
			},
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("diff (-got,+want:\n%s", diff)
	}
}
