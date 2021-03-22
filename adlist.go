package goebay

import (
	"context"
	"fmt"
	"net/http"
)

// AdListRepo represents an AdListRepo
// The AdListRepo resolves a list of ads
type AdListRepo struct {
	client *http.Client
}

// AdListResponse represents the response from the Fetch() call
type AdListResponse struct {
	Items []*AdListItem
}

// AdListItem a single item in the returned list
type AdListItem struct {
	ID              string
	Title           string
	Price           int
	PriceNegotiable bool
	Location        string
	ZipCode         string
	Link            string
}

// NewAdListRepo creates a new AdRepo, if client is nil, one will be created
func NewAdListRepo(client *http.Client) *AdListRepo {
	al := &AdListRepo{}

	if client == nil {
		al.client = &http.Client{}
	} else {
		al.client = client
	}

	return al
}

// Fetch fetches a list of ads based on the provided param
func (al *AdListRepo) Fetch(ctx context.Context, param *SearchParam) (*AdListResponse, error) {
	url := param.toURL()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := al.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	list, err := parseListHTML(resp.Body)

	if err != nil {
		return nil, err
	}

	return list, nil
}
