package portal

import (
	"fmt"

	"github.com/naag/terraform-provider-grafanacloud/internal/util"
)

type DataSourceList struct {
	Items []*DataSource
}
type DataSource struct {
	ID            int
	InstanceID    int
	InstanceSlug  string
	Name          string
	Type          string
	URL           string
	BasicAuth     int
	BasicAuthUser string
}

func (c *Client) ListDataSources(stack string) (*DataSourceList, error) {
	url := fmt.Sprintf("instances/%s/datasources", stack)
	resp, err := c.client.R().
		SetResult(&DataSourceList{}).
		Get(url)

	if err := util.HandleError(err, resp, "Failed to list Grafana data sources"); err != nil {
		return nil, err
	}

	return resp.Result().(*DataSourceList), nil
}

func (ds *DataSource) IsAlertmanager() bool {
	return ds.Type == "grafana-alertmanager-datasource"
}
