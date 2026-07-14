package discogs

import (
	"encoding/json"
	"fmt"
	"io"
)

// ReleaseDetail is the full response from GET /releases/{id}.
type ReleaseDetail struct {
	ID                int                    `json:"id"`
	Title             string                 `json:"title"`
	Year              int                    `json:"year"`
	Released          string                 `json:"released"`
	ReleasedFormatted string                 `json:"released_formatted"`
	Country           string                 `json:"country"`
	Status            string                 `json:"status"`
	Notes             string                 `json:"notes"`
	MasterID          int                    `json:"master_id"`
	MasterURL         string                 `json:"master_url"`
	ResourceURL       string                 `json:"resource_url"`
	URI               string                 `json:"uri"`
	Thumb             string                 `json:"thumb"`
	LowestPrice       float64                `json:"lowest_price"`
	NumForSale        int                    `json:"num_for_sale"`
	EstimatedWeight   int                    `json:"estimated_weight"`
	FormatQuantity    int                    `json:"format_quantity"`
	DataQuality       string                 `json:"data_quality"`
	Genres            []string               `json:"genres"`
	Styles            []string               `json:"styles"`
	Artists           []ReleaseArtist        `json:"artists"`
	ExtraArtists      []ReleaseArtist        `json:"extraartists"`
	Labels            []ReleaseLabel         `json:"labels"`
	Companies         []ReleaseLabel         `json:"companies"`
	Formats           []ReleaseFormat        `json:"formats"`
	Tracklist         []Track                `json:"tracklist"`
	Identifiers       []Identifier           `json:"identifiers"`
	Images            []Image                `json:"images"`
	Videos            []Video                `json:"videos"`
	Community         ReleaseDetailCommunity `json:"community"`
}

type ReleaseArtist struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ANV         string `json:"anv"`
	Join        string `json:"join"`
	Role        string `json:"role"`
	Tracks      string `json:"tracks"`
	ResourceURL string `json:"resource_url"`
}

type ReleaseLabel struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	CatNo          string `json:"catno"`
	EntityType     string `json:"entity_type"`
	EntityTypeName string `json:"entity_type_name"`
	ResourceURL    string `json:"resource_url"`
}

type ReleaseFormat struct {
	Name         string   `json:"name"`
	Qty          string   `json:"qty"`
	Descriptions []string `json:"descriptions"`
}

type Track struct {
	Position     string          `json:"position"`
	Title        string          `json:"title"`
	Duration     string          `json:"duration"`
	Type         string          `json:"type_"`
	Artists      []ReleaseArtist `json:"artists"`
	ExtraArtists []ReleaseArtist `json:"extraartists"`
}

type Identifier struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Image struct {
	Type        string `json:"type"`
	URI         string `json:"uri"`
	URI150      string `json:"uri150"`
	ResourceURL string `json:"resource_url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

type Video struct {
	URI         string `json:"uri"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Embed       bool   `json:"embed"`
}

type Community struct {
	Want int `json:"want"`
	Have int `json:"have"`
}

type ReleaseDetailCommunity struct {
	Have        int    `json:"have"`
	Want        int    `json:"want"`
	Status      string `json:"status"`
	DataQuality string `json:"data_quality"`
	Rating      Rating `json:"rating"`
}

type Rating struct {
	Average float64 `json:"average"`
	Count   int     `json:"count"`
}

// GetRelease fetches full details for a single release by its Discogs ID.
func (c *Client) GetRelease(id int) (*ReleaseDetail, error) {
	endpoint := fmt.Sprintf("%s/releases/%d", baseURL, id)

	req, err := c.newRequest(endpoint)
	if err != nil {
		return nil, fmt.Errorf("discogs: building request: %w", err)
	}

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

	var detail ReleaseDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, fmt.Errorf("discogs: decoding response: %w", err)
	}

	return &detail, nil
}
