package grafanacloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/portal"
)

func dataSourceStacks() *schema.Resource {
	s := baseStackSchema()
	s["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the stack to read (as slug name).",
	}

	return &schema.Resource{
		Description: "Reads all Grafana Cloud stacks which are provisioned inside the organisation.",
		ReadContext: dataSourceStacksRead,
		Schema: map[string]*schema.Schema{
			"stacks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: s,
				},
			},
		},
	}
}

func dataSourceStacksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	p := m.(*Provider)

	stacks, err := listStacks(p)
	if err != nil {
		return diag.FromErr(err)
	}

	schemaStacks := stackListToSchema(stacks)
	if err := d.Set("stacks", schemaStacks); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("all-stacks")

	return diags
}

func stackListToSchema(stacks []*portal.Stack) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	for _, stack := range stacks {
		result = append(result, map[string]interface{}{
			"id":                   stack.ID,
			"name":                 stack.Name,
			"prometheus_url":       stack.HmInstancePromURL,
			"prometheus_user_id":   stack.HmInstancePromID,
			"alertmanager_url":     stack.AmInstanceURL,
			"alertmanager_user_id": stack.AmInstanceID,
		})
	}

	return result
}

func listStacks(p *Provider) ([]*portal.Stack, error) {
	result := make([]*portal.Stack, 0)
	resp, err := p.Client.ListStacks(p.Organisation)
	if err != nil {
		return nil, err
	}

	for _, stack := range resp.Items {
		alertmanager, err := findAlertmanagerDatasource(p, stack)
		if err != nil {
			return nil, err
		}

		if alertmanager != nil {
			stack.AmInstanceURL = alertmanager.URL
		}

		result = append(result, stack)
	}

	return result, nil
}

func findAlertmanagerDatasource(p *Provider, stack *portal.Stack) (*portal.DataSource, error) {
	ds, err := p.Client.ListDataSources(stack.Slug)
	if err != nil {
		return nil, fmt.Errorf("error while locating Alertmanager instance for stack %s: %v", stack.Slug, err)
	}

	for _, datasource := range ds.Items {
		if datasource.IsAlertmanager() {
			return datasource, nil
		}
	}

	return nil, nil
}

func baseStackSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the Stack.",
		},
		"prometheus_url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Base URL of the Prometheus instance configured for this stack.",
		},
		"prometheus_user_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "User ID of the Prometheus instance configured for this stack.",
		},
		"alertmanager_url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Base URL of the Alertmanager instance configured for this stack. Please note that since this URL isn't provided by the Grafana Cloud API, this provider tries to obtain it from the Grafana data sources instead.",
		},
		"alertmanager_user_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "User ID of the Alertmanager instance configured for this stack.",
		},
	}
}
