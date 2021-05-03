package grafanacloud_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/naag/terraform-provider-grafanacloud/grafanacloud"
	"github.com/naag/terraform-provider-grafanacloud/internal/mock"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider
	grafanaCloudMock *mock.GrafanaCloud
)

const (
	EnvStack = "GRAFANA_CLOUD_STACK"
	EnvMock  = "GRAFANA_CLOUD_MOCK"
)

func TestMain(m *testing.M) {
	startMock()
	if grafanaCloudMock != nil {
		defer grafanaCloudMock.Close()
	}

	testAccProvider = grafanacloud.NewProvider("0.0.1")()
	testAccProviders = map[string]*schema.Provider{
		"grafanacloud": testAccProvider,
	}

	os.Exit(m.Run())
}

func startMock() {
	if os.Getenv(EnvMock) == "1" {
		orgName := os.Getenv(grafanacloud.EnvOrganisation)
		stackName := os.Getenv(EnvStack)

		grafanaCloudMock = mock.NewGrafanaCloud().
			Start().
			WithOrganisation(orgName).
			WithStack(stackName, orgName)

		os.Setenv(grafanacloud.EnvURL, fmt.Sprintf("%s/api", grafanaCloudMock.Server.URL))
	}
}
