package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const apiBase = "https://api.github.com"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{httpClient: &http.Client{}}
}

type Key struct {
	Key string `json:"key"`
}

func (c *Client) UserKeys(user string) ([]Key, error) {
	req, err := http.NewRequest("GET", apiBase+"/users/"+user+"/keys", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gitbastion")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var keys []Key
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, err
	}
	return keys, nil
}
