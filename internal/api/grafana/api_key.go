package grafana

import (
	"fmt"
	"strconv"
	"time"

	"github.com/naag/terraform-provider-grafanacloud/internal/util"
	"github.com/relvacode/iso8601"
)

type CreateAPIKey struct {
	Name          string `json:"name"`
	Role          string `json:"role"`
	SecondsToLive int    `json:"secondsToLive"`
}

type APIKeyList struct {
	Keys []*APIKey
}

type APIKey struct {
	ID         int
	Name       string
	Key        string
	Expiration string
	Role       string
}

func (c *Client) ListAPIKeys(includeExpired bool) (*APIKeyList, error) {
	var apiKeys []*APIKey
	url := "api/auth/keys"

	resp, err := c.client.R().
		SetResult(&apiKeys).
		SetQueryParam("includeExpired", strconv.FormatBool(includeExpired)).
		Get(url)

	if err := util.HandleError(err, resp, "Failed to list Grafana API keys"); err != nil {
		return nil, err
	}

	return &APIKeyList{
		Keys: apiKeys,
	}, nil
}

func (c *Client) DeleteAPIKey(id int) error {
	url := fmt.Sprintf("api/auth/keys/%d", id)

	resp, err := c.client.R().
		Delete(url)

	if err := util.HandleError(err, resp, "Failed to delete Grafana API key"); err != nil {
		return err
	}

	return nil
}

func (k *APIKey) IsExpired() (bool, error) {
	expires, err := iso8601.ParseString(k.Expiration)
	if err != nil {
		return false, err
	}

	now := time.Now()
	return now.After(expires), nil
}

func (l *APIKeyList) AddKey(k *APIKey) {
	l.Keys = append(l.Keys, k)
}

func (l *APIKeyList) FindByID(id int) *APIKey {
	for _, k := range l.Keys {
		if k.ID == id {
			return k
		}
	}

	return nil
}

func (l *APIKeyList) DeleteByID(id int) {
	newKeys := make([]*APIKey, 0)

	for _, k := range l.Keys {
		if k.ID != id {
			newKeys = append(newKeys, k)
		}
	}

	l.Keys = newKeys
}
