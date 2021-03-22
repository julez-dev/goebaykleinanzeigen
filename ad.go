package goebay

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type AdRepo struct {
	client *http.Client
}

type AdItem struct {
	ID              string
	ListedSince     time.Time
	Title           string
	Price           int
	PriceNegotiable bool
	Location        string
	ZipCode         string
	Link            string
	Description     string
	Details         map[string]string
	Extras          []string
	Seller          *Seller
}

type Seller struct {
	Name         string
	ActiveSince  time.Time
	Friendliness string
	Rating       string
}

func NewAdRepo(client *http.Client) *AdRepo {
	ar := &AdRepo{}

	if client == nil {
		ar.client = &http.Client{}
	} else {
		ar.client = client
	}

	return ar
}

func (ar *AdRepo) Fetch(ctx context.Context, url string) (*AdItem, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := ar.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code or empty content length")
	}

	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(body))
	// r := bytes.NewReader(body)
	// item, err := parseAdHTML(r)

	item, err := parseAdHTML(resp.Body)

	if err != nil {
		return nil, err
	}

	return item, nil
}
