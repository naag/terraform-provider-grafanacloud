package grafanacloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/portal"
)

func resourceStack() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a single Stack in Grafana Cloud.",
		CreateContext: resourceStackCreate,
		ReadContext:   resourceStackRead,
		DeleteContext: resourceStackDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Grafana Cloud stack.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Grafana Cloud stack.",
			},
			"slug": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Slug name of the Grafana Cloud stack.",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Custom URL for the Grafana instance. Must have a CNAME setup to point to `.grafana.net` before creating the stack.",
			},
		},
	}
}

func resourceStackCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	p := m.(*Provider)

	req := &portal.CreateStack{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
		URL:  d.Get("url").(string),
	}

	resp, err := p.Client.CreateStack(req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(resp.ID))

	return resourceStackRead(ctx, d, m)
}

func resourceStackRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	p := m.(*Provider)

	slug := d.Get("slug").(string)
	resp, err := p.Client.GetStack(p.Organisation, slug)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("slug", resp.Slug); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url", resp.URL); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceStackDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	p := m.(*Provider)

	slug := d.Get("slug").(string)
	err := p.Client.DeleteStack(slug)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
