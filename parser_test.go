package goebaykleinanzeigen

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func timeCanPanic(t *testing.T, text string) time.Time {
	time, err := time.Parse("02.01.2006", text)

	if err != nil {
		t.Errorf("helper: could not parse time")
	}

	return time
}

func Test_parsePrice(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		price      int
		negotiable bool
	}{
		{
			name:       "only-negotiable",
			text:       "VB",
			price:      0,
			negotiable: true,
		},
		{
			name:       "only-price",
			text:       "4.500 €",
			price:      4500,
			negotiable: false,
		},
		{
			name:       "price-and-negotiable",
			text:       "4.500 € VB",
			price:      4500,
			negotiable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, negotiable, _ := parsePrice(tt.text)
			if price != tt.price {
				t.Errorf("parsePrice() price = %v, want %v", price, tt.price)
			}
			if negotiable != tt.negotiable {
				t.Errorf("parsePrice() negotiable = %v, want %v", negotiable, tt.negotiable)
			}
		})
	}
}

func Test_parseLocation(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		zip      string
		location string
	}{
		{
			name:     "zip-and-location",
			text:     "10119 Mitte\t",
			zip:      "10119",
			location: "Mitte",
		},
		{
			name:     "zip-and-location-multiple",
			text:     "16540 Hohen Neuendorf",
			zip:      "16540",
			location: "Hohen Neuendorf",
		},
		{
			name:     "zip-and-location-multiple-distance",
			text:     "16540 Hohen Neuendorf\n(13km)",
			zip:      "16540",
			location: "Hohen Neuendorf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zip, location := parseLocation(tt.text)
			if zip != tt.zip {
				t.Errorf("parseLocation() zip = %v, want %v", zip, tt.zip)
			}
			if location != tt.location {
				t.Errorf("parseLocation() location = %v, want %v", location, tt.location)
			}
		})
	}
}

func Test_parseExtraInfo(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		listedSince time.Time
		id          string
	}{
		{
			name:        "normal",
			text:        "22.03.2021\n Anzeigennr.: 1708911891",
			listedSince: timeCanPanic(t, "22.03.2021"),
			id:          "1708911891",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listedSince, id := parseExtraInfo(tt.text)
			if listedSince != tt.listedSince {
				t.Errorf("parseExtraInfo() listedSince = %v, want %v", listedSince, tt.listedSince)
			}
			if id != tt.id {
				t.Errorf("parseExtraInfo() id = %v, want %v", id, tt.id)
			}
		})
	}
}

func Test_parseDetail(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		key   string
		value string
	}{
		{
			name:  "km",
			text:  "Kilometerstand\n38.600 km",
			key:   "Kilometerstand",
			value: "38.600 km",
		},
		{
			name:  "doors",
			text:  "Anzahl Türen\n2/3",
			key:   "Anzahl Türen",
			value: "2/3",
		},
		{
			name:  "spaces-value",
			text:  "Fahrzeugzustand\nUnbeschädigtes Fahrzeug",
			key:   "Fahrzeugzustand",
			value: "Unbeschädigtes Fahrzeug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, value := parseDetail(tt.text)
			if key != tt.key {
				t.Errorf("parseDetail() key = %v, want %v", key, tt.key)
			}
			if value != tt.value {
				t.Errorf("parseDetail() value = %v, want %v", value, tt.value)
			}
		})
	}
}

func Test_parseRating(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		rating string
	}{
		{
			name:   "empty",
			text:   "",
			rating: "",
		},
		{
			name:   "top",
			text:   "Zufriedenheit: TOP",
			rating: "TOP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRating(tt.text); got != tt.rating {
				t.Errorf("parseRating() = %v, want %v", got, tt.rating)
			}
		})
	}
}

func Test_parseActiveSince(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		activeSince time.Time
	}{
		{
			name:        "normal",
			text:        "Aktiv seit 06.10.2012",
			activeSince: timeCanPanic(t, "06.10.2012"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			activeSince, err := parseActiveSince(tt.text)

			if err != nil {
				t.Errorf("did not expect error, got %w", err)
			}

			if activeSince != tt.activeSince {
				t.Errorf("parseActiveSince() activeSince = %v, want %v", activeSince, tt.activeSince)
			}
		})
	}
}

func Test_parseAdHTML(t *testing.T) {
	files, err := filepath.Glob("testdata/aditem/*.html")

	if err != nil {
		t.Fatal(err)
	}

	for _, htmlFilePath := range files {
		htmlFilename := filepath.Base(htmlFilePath)
		jsonFilePath := "testdata/aditem/" + strings.Replace(htmlFilename, ".html", ".json", 1)

		t.Run(htmlFilename, func(t *testing.T) {
			// html, err := os.ReadFile(htmlFilePath)
			html, err := ioutil.ReadFile(htmlFilePath)

			if err != nil {
				t.Fatal(err)
			}

			// jsonRaw, err := os.ReadFile(jsonFilePath)
			jsonRaw, err := ioutil.ReadFile(jsonFilePath)

			if err != nil {
				t.Fatal(err)
			}

			adItem := &AdItem{}
			_ = json.Unmarshal(jsonRaw, adItem)

			returnedItem, err := parseAdHTML(bytes.NewReader(html))

			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(returnedItem, adItem) {
				t.Errorf("parseAdHTML() got = %v, want %v", returnedItem, adItem)
			}
		})
	}
}
