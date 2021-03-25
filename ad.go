package goebay

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// AdRepo represents an AdRepo
// The AdRepo resolves a single ad
type AdRepo struct {
	client *http.Client
}

// AdItem represents a single AdItem
type AdItem struct {
	ID              string    `json:"id"`
	ListedSince     time.Time `json:"listed_since"`
	Title           string    `json:"title"`
	Price           int       `json:"price"`
	PriceNegotiable bool      `json:"price_negotiable"`
	Location        string    `json:"location"`
	ZipCode         string    `json:"zip_code"`
	Link            string    `json:"link"`
	Description     string    `json:"description"`
	Details         []*Detail `json:"details"`
	Extras          []string  `json:"extras"`
	Seller          *Seller   `json:"seller"`
}

// Seller represents the seller of an aditem
type Seller struct {
	Name         string    `json:"name"`
	ActiveSince  time.Time `json:"active_since"`
	Friendliness string    `json:"friendliness"`
	Rating       string    `json:"rating"`
}

type Detail struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// NewAdRepo creates a new AdRepo, if client is nil, one will be created
func NewAdRepo(client *http.Client) *AdRepo {
	ar := &AdRepo{}

	if client == nil {
		ar.client = &http.Client{}
	} else {
		ar.client = client
	}

	return ar
}

// Fetch fetches a single AdItem
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
