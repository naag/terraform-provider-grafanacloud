package grafanacloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/portal"
)

func dataSourceStack() *schema.Resource {
	s := baseStackSchema()
	s["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the stack to read (as slug name).",
	}

	return &schema.Resource{
		Description: "Reads a single Grafana Cloud stack from the organisation by the given name.",
		ReadContext: dataSourceStackRead,
		Schema:      s,
	}
}
func dataSourceStackRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	p := m.(*Provider)
	nameFilter := d.Get("name").(string)

	stacks, err := listStacks(p)
	if err != nil {
		return diag.FromErr(err)
	}

	stacks = filterStacksByName(stacks, nameFilter)
	if len(stacks) != 1 {
		return diag.Errorf("Expected to find a single stack, found %d stacks. Please check the name attribute", len(stacks))
	}

	d.SetId(fmt.Sprint(stacks[0].ID))

	if err := d.Set("name", stacks[0].Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("slug", stacks[0].Slug); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("prometheus_url", stacks[0].HmInstancePromURL); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("prometheus_user_id", stacks[0].HmInstancePromID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("alertmanager_url", stacks[0].AmInstanceURL); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("alertmanager_user_id", stacks[0].AmInstanceID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func filterStacksByName(stacks []*portal.Stack, nameFilter string) []*portal.Stack {
	result := make([]*portal.Stack, 0)

	for _, stack := range stacks {
		if nameFilter == "" || stack.Slug == nameFilter {
			result = append(result, stack)
		}
	}

	return result
}
