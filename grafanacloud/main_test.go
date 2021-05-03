package grafanacloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/naag/terraform-provider-grafanacloud/internal/mock"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider
	grafanaCloudMock *mock.GrafanaCloud
)

func TestMain(m *testing.M) {
	startMock()
	if grafanaCloudMock != nil {
		defer grafanaCloudMock.Close()
	}

	testAccProvider = Provider("0.0.1")()
	testAccProviders = map[string]*schema.Provider{
		"grafanacloud": testAccProvider,
	}

	os.Exit(m.Run())
}

func startMock() {
	if os.Getenv("GRAFANA_CLOUD_MOCK") == "1" {
		orgName := os.Getenv(EnvOrganisation)
		stackName := os.Getenv("GRAFANA_CLOUD_STACK")

		grafanaCloudMock = mock.NewGrafanaCloud().
			Start().
			WithOrganisation(orgName).
			WithStack(stackName, orgName)

		os.Setenv(EnvURL, fmt.Sprintf("%s/api", grafanaCloudMock.Server.URL))
	}
}
