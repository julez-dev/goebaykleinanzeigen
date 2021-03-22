package goebay

import (
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func parseListHTML(body io.Reader) (*AdListResponse, error) {
	doc, err := goquery.NewDocumentFromReader(body)

	if err != nil {
		return nil, err
	}

	selector := doc.Find(".aditem")
	listItems := make([]*AdListItem, 0, len(selector.Nodes))

	response := &AdListResponse{
		Items: listItems,
	}

	selector.Each(func(_ int, s *goquery.Selection) {
		listItem := &AdListItem{}

		if id, ok := s.Attr("data-adid"); ok {
			listItem.ID = id
			listItem.Link = BaseURL + "/s-anzeige/" + listItem.ID
		}

		listItem.Title = s.Find(".ellipsis").First().Text()

		priceText := strings.TrimSpace(s.Find(".aditem-main--middle--price").First().Text())
		price, negotiable, _ := parsePrice(priceText)

		listItem.PriceNegotiable = negotiable
		listItem.Price = price

		locationText := strings.TrimSpace(s.Find(".aditem-main--top--left").First().Text())
		zip, location := parseLocation(locationText)

		listItem.Location = location
		listItem.ZipCode = zip

		response.Items = append(response.Items, listItem)
	})

	return response, nil
}

func parsePrice(text string) (int, bool, error) {

	splits := strings.Split(text, " ")

	if len(splits) > 0 {
		splits[0] = strings.ReplaceAll(splits[0], ".", "")

		if splits[0] == "VB" {
			return 0, true, nil
		}

		if len(splits) == 2 || len(splits) == 3 {
			price, err := strconv.ParseInt(splits[0], 10, 32)

			if err != nil {
				return 0, false, err
			}

			return int(price), false, nil
		}

	}

	return 0, false, nil
}

func parseLocation(text string) (string, string) {
	text = strings.TrimSpace(text)
	splits := strings.SplitN(text, " ", 2)

	if len(splits) == 2 {
		// the location may also include the distance
		// if a radius is provided in the parameters.
		// In this case cut location until the newline
		newLineAt := strings.Index(splits[1], "\n")

		if newLineAt >= 0 {
			return splits[0], splits[1][:newLineAt]
		}

		return splits[0], splits[1]
	}

	return "", ""
}

func parseAdHTML(body io.Reader) (*AdItem, error) {
	doc, err := goquery.NewDocumentFromReader(body)

	if err != nil {
		return nil, err
	}

	ad := &AdItem{}
	seller := &Seller{}
	ad.Seller = seller

	ad.Title = strings.TrimSpace(doc.Find("#viewad-title").First().Text())

	priceText := strings.TrimSpace(doc.Find("#viewad-price").First().Text())
	price, negotiable, _ := parsePrice(priceText)

	ad.PriceNegotiable = negotiable
	ad.Price = price

	locationText := strings.TrimSpace(doc.Find("#viewad-locality").First().Text())
	zip, location := parseLocation(locationText)

	ad.Location = location
	ad.ZipCode = zip

	// It's not directly possible to scrape the number of views from the html directly since its injected by javascript.
	// The current number of views are available under https://www.ebay-kleinanzeigen.de/s-vac-inc-get.json?adId={ID}
	// fmt.Println(doc.Find(".textcounter").Html())
	date, id := parseExtraInfo(strings.TrimSpace(doc.Find("#viewad-extra-info").First().Text()))

	ad.ListedSince = date
	ad.ID = id
	ad.Link = BaseURL + "/s-anzeige/" + id

	detailsSelector := doc.Find(".addetailslist--detail")
	ad.Details = make(map[string]string, len(detailsSelector.Nodes))

	detailsSelector.Each(func(_ int, s *goquery.Selection) {
		key, val := parseDetail(strings.TrimSpace(s.Text()))
		ad.Details[key] = val
	})

	extrasSelector := doc.Find(".checktag")
	ad.Extras = make([]string, 0, len(extrasSelector.Nodes))

	extrasSelector.Each(func(_ int, s *goquery.Selection) {
		ad.Extras = append(ad.Extras, s.Text())
	})

	seller.Name = strings.TrimSpace(doc.Find(".text-bold.text-bigger.text-force-linebreak").First().Text())
	seller.Rating = parseRating(strings.TrimSpace(doc.Find(".userbadges-vip.userbadges-profile-rating").First().Text()))
	seller.Friendliness = strings.TrimSpace(doc.Find(".userbadges-vip.userbadges-profile-friendliness").First().Text())

	time, err := parseActiveSince(strings.TrimSpace(doc.Find(".text-light.text-light-seller-info").First().Text()))

	// We don't really care if this returns an error
	// just don't set the field in this case
	if err == nil {
		seller.ActiveSince = time
	}

	ad.Description = strings.TrimSpace(doc.Find("#viewad-description-text").First().Text())

	return ad, nil
}

func parseExtraInfo(text string) (time.Time, string) {
	text = strings.ReplaceAll(text, "\n", "")
	splits := strings.SplitN(text, " ", 2)

	if len(splits) > 1 {
		splits[1] = strings.TrimSpace(splits[1])
		splits[1] = strings.TrimPrefix(splits[1], "Anzeigennr.: ")

		date, _ := time.Parse("02.01.2006", splits[0])

		return date, splits[1]
	}

	return time.Time{}, ""
}

func parseDetail(text string) (string, string) {
	splits := strings.SplitN(text, "\n", 2)
	return strings.TrimSpace(splits[0]), strings.TrimSpace(splits[1])
}

func parseRating(text string) string {
	return strings.TrimPrefix(text, "Zufriedenheit: ")
}

func parseActiveSince(text string) (time.Time, error) {
	if text == "" {
		return time.Time{}, nil
	}

	strip := "Aktiv seit "
	index := strings.Index(text, strip)
	text = text[index+len(strip):]
	splits := strings.Split(text, "\n")

	return time.Parse("02.01.2006", splits[0])
}
