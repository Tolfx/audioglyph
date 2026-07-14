package discogs

import (
	"log/slog"
	"net/http"
	"time"
)

const baseURL = "https://api.discogs.com"
const userAgent = "audioglyph/0.1 +https://github.com/Tolfx/audioglyph"
const statusOK = http.StatusOK

type Client struct {
	logger     *slog.Logger
	httpClient *http.Client
	token      string
}

func NewClient(logger *slog.Logger, token string) *Client {
	return &Client{
		logger:     logger,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		token:      token,
	}
}

// newRequest creates a GET request with standard headers applied.
func (c *Client) newRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	if c.token != "" {
		req.Header.Set("Authorization", "Discogs token="+c.token)
	}
	return req, nil
}
