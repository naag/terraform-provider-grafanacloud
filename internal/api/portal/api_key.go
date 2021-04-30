package portal

import (
	"fmt"

	"github.com/naag/terraform-provider-grafanacloud/internal/api/grafana"
	"github.com/naag/terraform-provider-grafanacloud/internal/util"
)

type CreateAPIKey struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type APIKeyList struct {
	Items []*APIKey
}

type APIKey struct {
	ID    int
	Name  string
	Role  string
	Token string
}

// This function creates a API key inside the Grafana instance running in stack `stack`. It's used in order
// to provision API keys inside Grafana while just having access to a Grafana Cloud API key.
//
// Plese note that this is a beta feature and might change in the future.
//
// See https://grafana.com/docs/grafana-cloud/api/#create-grafana-api-keys for more information.
func (c *Client) CreateGrafanaAPIKey(r *grafana.CreateAPIKey, stack string) (*grafana.APIKey, error) {
	url := fmt.Sprintf("instances/%s/api/auth/keys", stack)
	resp, err := c.client.R().
		SetBody(r).
		SetResult(&grafana.APIKey{}).
		Post(url)

	if err := util.HandleError(err, resp, "Failed to create Grafana API key through Grafana Cloud proxy route"); err != nil {
		return nil, err
	}

	return resp.Result().(*grafana.APIKey), nil
}

func (c *Client) CreateAPIKey(r *CreateAPIKey, org string) (*APIKey, error) {
	url := fmt.Sprintf("orgs/%s/api-keys", org)
	resp, err := c.client.R().
		SetBody(r).
		SetResult(&APIKey{}).
		Post(url)

	if err := util.HandleError(err, resp, "Failed to create Grafana Cloud Portal API key"); err != nil {
		return nil, err
	}

	return resp.Result().(*APIKey), nil
}

func (c *Client) ListAPIKeys(org string) (*APIKeyList, error) {
	url := fmt.Sprintf("orgs/%s/api-keys", org)
	resp, err := c.client.R().
		SetResult(&APIKeyList{}).
		Get(url)

	if err := util.HandleError(err, resp, "Failed to read Grafana Cloud Portal API key"); err != nil {
		return nil, err
	}

	return resp.Result().(*APIKeyList), nil
}

func (c *Client) DeleteAPIKey(org string, keyName string) error {
	url := fmt.Sprintf("orgs/%s/api-keys/%s", org, keyName)
	resp, err := c.client.R().
		Delete(url)

	if err := util.HandleError(err, resp, "Failed to delete Grafana Cloud Portal API key"); err != nil {
		return err
	}

	return nil
}

func (l *APIKeyList) FindByName(name string) *APIKey {
	for _, k := range l.Items {
		if k.Name == name {
			return k
		}
	}

	return nil
}
