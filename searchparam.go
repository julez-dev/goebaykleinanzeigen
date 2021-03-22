package goebay

import (
	"strconv"
	"strings"
)

type Category string
type Provider string
type OfferType string
type CountryID string
type Radius string

type SpecificParameter map[ParamName]string

const (
	BaseURL string = "https://www.ebay-kleinanzeigen.de"
)

const (
	Cars Category = "216"
)

const (
	Private    Provider = "privat"
	Commercial Provider = "gewerblich"
)

const (
	Offer  OfferType = "angebote"
	Search OfferType = "gesuche"
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

type SearchParam struct {
	Category          Category
	Provider          Provider
	OfferType         OfferType
	CountryID         CountryID
	Radius            Radius
	SpecificParameter SpecificParameter
	Page              int
	PriceFrom         int
	PriceTo           int
}

func (sp *SearchParam) fmtCategory() string {
	if sp.Category != "" {
		return "/c" + string(sp.Category)
	}

	return ""
}

func (sp *SearchParam) fmtPrice() string {
	if sp.PriceTo > 0 {
		return "/preis:" + strconv.Itoa(sp.PriceFrom) + ":" + strconv.Itoa(sp.PriceTo)
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

func (sp *SearchParam) fmtCountyID() string {
	if sp.CountryID != "" {
		return "l" + string(sp.CountryID)
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

func (sp *SearchParam) ToURL() string {
	sb := strings.Builder{}

	_, _ = sb.WriteString(sp.fmtOfferType())
	_, _ = sb.WriteString(sp.fmtPrice())
	_, _ = sb.WriteString(sp.fmtProvider())
	_, _ = sb.WriteString(sp.fmtPage())
	_, _ = sb.WriteString(sp.fmtCategory())
	_, _ = sb.WriteString(sp.fmtCountyID())
	_, _ = sb.WriteString(sp.fmtRadius())
	_, _ = sb.WriteString(sp.fmtSpecificParameter())

	params := strings.Trim(sb.String(), "/")

	return BaseURL + "/" + params
}
