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

func (l *StackList) FindBySlug(slug string) *Stack {
	for _, stack := range l.Items {
		if stack.Slug == slug {
			return stack
		}
	}

	return nil
}
