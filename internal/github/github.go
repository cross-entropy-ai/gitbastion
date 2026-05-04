package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const apiBase = "https://api.github.com"

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{},
	}
}

type Key struct {
	Key string `json:"key"`
}

type member struct {
	Login string `json:"login"`
}

func (c *Client) fetch(path string, result any) error {
	req, err := http.NewRequest("GET", apiBase+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "gitbastion")
	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) UserKeys(user string) ([]Key, error) {
	var keys []Key
	if err := c.fetch("/users/"+user+"/keys", &keys); err != nil {
		return nil, fmt.Errorf("fetch keys for %s: %w", user, err)
	}
	return keys, nil
}

// TeamMembers returns logins of all members in an org team.
// team is in "org/team-slug" format. Requires a token with read:org scope.
func (c *Client) TeamMembers(org, slug string) ([]string, error) {
	var logins []string
	for page := 1; ; page++ {
		var batch []member
		path := fmt.Sprintf("/orgs/%s/teams/%s/members?per_page=100&page=%d", org, slug, page)
		if err := c.fetch(path, &batch); err != nil {
			return nil, fmt.Errorf("fetch team %s/%s: %w", org, slug, err)
		}
		if len(batch) == 0 {
			break
		}
		for _, m := range batch {
			logins = append(logins, m.Login)
		}
	}
	return logins, nil
}
