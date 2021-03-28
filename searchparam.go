package goebaykleinanzeigen

import (
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

// Category represents a specific category id
type Category string

// Provider represents a specific provider (private or commercial seller)
type Provider string

// OfferType represents the offer type
type OfferType string

// LocationID represents a specifc location
// retrived from https://www.ebay-kleinanzeigen.de/s-ort-empfehlungen.json?query=Berlin for example
type LocationID string

// Radius represents the distance from the center of the location
type Radius string

// SpecificParameter represent specific parameters for each category
type SpecificParameter map[ParamName]string

const (
	baseURL string = "https://www.ebay-kleinanzeigen.de"
)

const (
	Family         Category = "17"
	Music          Category = "73"
	HouseAndGarden Category = "80"
	Jobs           Category = "102"
	Pets           Category = "130"
	Fashion        Category = "153"
	Electronics    Category = "161"
	Hobby          Category = "185"
	Property       Category = "195"
	Boats          Category = "211"
	Cars           Category = "216"
	CarParts       Category = "223"
	Tickets        Category = "231"
	Education      Category = "235"
	Services       Category = "297"
)

const (
	Private    Provider = "privat"
	Commercial Provider = "gewerblich"
)

const (
	Offer  OfferType = "angebote"
	Wanted OfferType = "gesuche"
)

const (
	WholePlace        Radius = ""
	FiveKM            Radius = "5"
	TenKM             Radius = "10"
	TwentyKM          Radius = "20"
	ThirtyKM          Radius = "30"
	FiftyKM           Radius = "50"
	OneHundredKM      Radius = "100"
	OneHundredFiftyKM Radius = "150"
	TwoHundredKM      Radius = "200"
)

// SearchParam holds all parameters for a search
type SearchParam struct {
	Category          Category          `json:"category"`
	Provider          Provider          `json:"provider"`
	OfferType         OfferType         `json:"offer_type"`
	Location          LocationID        `json:"location"`
	Radius            Radius            `json:"radius"`
	SpecificParameter SpecificParameter `json:"specific_parameter"`
	// Which page should be scraped
	Page int `json:"page"`
	// Price from in euro
	PriceFrom int `json:"price_from"`
	// Price to in euro
	PriceTo int `json:"price_to"`
}

// ParamsFromJSON creates a SearchParam struct from JSON
func ParamsFromJSON(data io.Reader) (*SearchParam, error) {
	params := &SearchParam{}

	err := json.NewDecoder(data).Decode(params)

	if err != nil && err != io.EOF {
		return nil, err
	}

	return params, nil
}

func (sp *SearchParam) ToJSON(out io.Writer) error {
	err := json.NewEncoder(out).Encode(sp)

	if err != nil {
		return err
	}

	return nil
}

func (sp *SearchParam) fmtCategory() string {
	if sp.Category != "" {
		return "/c" + string(sp.Category)
	}

	return ""
}

func (sp *SearchParam) fmtPrice() string {
	if sp.PriceTo > 0 {
		if sp.PriceFrom > 0 {
			return "/preis:" + strconv.Itoa(sp.PriceFrom) + ":" + strconv.Itoa(sp.PriceTo)
		}

		return "/preis::" + strconv.Itoa(sp.PriceTo)
	}

	if sp.PriceFrom > 0 {
		return "/preis:" + strconv.Itoa(sp.PriceFrom) + ":"
	}

	return ""
}

func (sp *SearchParam) fmtPage() string {
	if sp.Page > 0 {
		return "/seite:" + strconv.Itoa(sp.Page)
	}

	return ""
}

func (sp *SearchParam) fmtProvider() string {
	providerString := string(sp.Provider)

	if providerString != "" {
		return "/anbieter:" + providerString
	}

	return ""
}

func (sp *SearchParam) fmtOfferType() string {
	offerString := string(sp.OfferType)

	if offerString != "" {
		return "/anzeige:" + offerString
	}

	return ""
}

func (sp *SearchParam) fmtLocationID() string {
	if sp.Location != "" {
		return "l" + string(sp.Location)
	}

	return ""
}

func (sp *SearchParam) fmtRadius() string {
	if sp.Radius != "" {
		return "r" + string(sp.Radius)
	}

	return ""
}

func (sp *SearchParam) fmtSpecificParameter() string {
	sb := strings.Builder{}

	if sp.SpecificParameter != nil {
		for k, v := range sp.SpecificParameter {
			_, _ = sb.WriteString("+" + string(k) + ":" + v)
		}
	}

	return sb.String()
}

func (sp *SearchParam) toURL() string {
	sb := strings.Builder{}

	if sp.Page < 1 {
		sp.Page = 1
	}

	_, _ = sb.WriteString(sp.fmtOfferType())
	_, _ = sb.WriteString(sp.fmtPrice())
	_, _ = sb.WriteString(sp.fmtProvider())
	_, _ = sb.WriteString(sp.fmtPage())
	_, _ = sb.WriteString(sp.fmtCategory())
	_, _ = sb.WriteString(sp.fmtLocationID())
	_, _ = sb.WriteString(sp.fmtRadius())
	_, _ = sb.WriteString(sp.fmtSpecificParameter())

	params := strings.Trim(sb.String(), "/")

	return baseURL + "/" + params
}
