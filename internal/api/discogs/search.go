package discogs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

// SearchResult is the top-level response from the Discogs search endpoint.
type SearchResult struct {
	Pagination Pagination `json:"pagination"`
	Results    []Release  `json:"results"`
}

type Pagination struct {
	PerPage int            `json:"per_page"`
	Pages   int            `json:"pages"`
	Page    int            `json:"page"`
	Items   int            `json:"items"`
	URLs    PaginationURLs `json:"urls"`
}

type PaginationURLs struct {
	Last string `json:"last"`
	Next string `json:"next"`
	Prev string `json:"prev"`
}

// Release is a search result entry.
type Release struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Year        string    `json:"year"`
	Country     string    `json:"country"`
	Format      []string  `json:"format"`
	Genre       []string  `json:"genre"`
	Style       []string  `json:"style"`
	Label       []string  `json:"label"`
	CatNo       string    `json:"catno"`
	Barcode     []string  `json:"barcode"`
	URI         string    `json:"uri"`
	ResourceURL string    `json:"resource_url"`
	Thumb       string    `json:"thumb"`
	Type        string    `json:"type"`
	Community   Community `json:"community"`
}

// SearchByBarcode searches Discogs for releases matching the given barcode.
func (c *Client) SearchByBarcode(barcode string) (*SearchResult, error) {
	endpoint := fmt.Sprintf("%s/database/search", baseURL)

	req, err := c.newRequest(endpoint)
	if err != nil {
		return nil, fmt.Errorf("discogs: building request: %w", err)
	}

	q := url.Values{}
	q.Set("barcode", barcode)
	q.Set("type", "release")
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("discogs: request failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.logger.Error("discogs: failed to close response body", "err", err)
		}
	}(resp.Body)

	if resp.StatusCode != statusOK {
		return nil, fmt.Errorf("discogs: unexpected status %s", resp.Status)
	}

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("discogs: decoding response: %w", err)
	}

	return &result, nil
}
