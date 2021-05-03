package portal

import (
	"fmt"

	"github.com/naag/terraform-provider-grafanacloud/internal/util"
)

type StackList struct {
	Items []*Stack
}

type Stack struct {
	ID                   int
	OrgID                int
	OrgSlug              string
	OrgName              string
	Name                 string
	URL                  string
	Status               string
	Slug                 string
	HmInstancePromID     int
	HmInstancePromURL    string
	HmInstancePromStatus string
	AmInstanceID         int
	AmInstanceURL        string
}

type CreateStack struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	URL  string `json:"url"`
}

func (c *Client) CreateStack(r *CreateStack) (*Stack, error) {
	url := "instances"
	resp, err := c.client.R().
		SetBody(r).
		SetResult(&Stack{}).
		Post(url)

	if err := util.HandleError(err, resp, "Failed to create Grafana Cloud stack"); err != nil {
		return nil, err
	}

	return resp.Result().(*Stack), nil
}

func (c *Client) ListStacks(org string) (*StackList, error) {
	url := fmt.Sprintf("orgs/%s/instances", org)
	resp, err := c.client.R().
		SetResult(&StackList{}).
		Get(url)

	if err := util.HandleError(err, resp, "Failed to list Grafana Cloud stacks"); err != nil {
		return nil, err
	}

	return resp.Result().(*StackList), nil
}

func (c *Client) GetStack(org, stackSlug string) (*Stack, error) {
	stacks, err := c.ListStacks(org)
	if err != nil {
		return nil, err
	}

	stack := stacks.FindBySlug(stackSlug)
	return stack, nil
}

func (c *Client) DeleteStack(stackSlug string) error {
	url := fmt.Sprintf("instances/%s", stackSlug)
	resp, err := c.client.R().
		Delete(url)

	if err := util.HandleError(err, resp, "Failed to delete Grafana Cloud stack"); err != nil {
		return err
	}

	return nil
}

func (l *StackList) FindBySlug(slug string) *Stack {
	for _, stack := range l.Items {
		if stack.Slug == slug {
			return stack
		}
	}

	return nil
}
