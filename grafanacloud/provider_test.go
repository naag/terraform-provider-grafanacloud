package grafanacloud

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/naag/terraform-provider-grafanacloud/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	if err := Provider("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderConfigure(t *testing.T) {
	mock := mock.NewGrafanaCloud().Start()
	defer mock.Close()

	resourceSchema := map[string]*schema.Schema{
		"url": {
			Type: schema.TypeString,
		},
		"api_key": {
			Type: schema.TypeString,
		},
		"organisation": {
			Type: schema.TypeString,
		},
	}

	resourceDataMap := map[string]interface{}{
		"url":          os.Getenv(EnvURL),
		"api_key":      os.Getenv(EnvAPIKey),
		"organisation": os.Getenv(EnvOrganisation),
	}
	resourceLocalData := schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)

	configureFunc := configureProvider("0.0.1", &schema.Provider{TerraformVersion: "0.15"})
	provider, err := configureFunc(context.TODO(), resourceLocalData)
	require.Nil(t, err)

	_, ok := provider.(*grafanaCloudProvider)
	require.True(t, ok)
}

func getProvider(p *schema.Provider) *grafanaCloudProvider {
	return p.Meta().(*grafanaCloudProvider)
}
