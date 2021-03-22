package goebay

import (
	"context"
	"fmt"
	"net/http"
)

type AdListRepo struct {
	client *http.Client
}

type AdListResponse struct {
	Items []*AdListItem
}

type AdListItem struct {
	ID              string
	Title           string
	Price           int
	PriceNegotiable bool
	Location        string
	ZipCode         string
	Link            string
}

func NewAdListRepo(client *http.Client) *AdListRepo {
	al := &AdListRepo{}

	if client == nil {
		al.client = &http.Client{}
	} else {
		al.client = client
	}

	return al
}

func (al *AdListRepo) Fetch(ctx context.Context, param *SearchParam) (*AdListResponse, error) {
	url := param.ToURL()
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
