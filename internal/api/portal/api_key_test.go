package portal_test

/*
import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/naag/terraform-provider-grafanacloud/internal/api/grafana"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/portal"
	"github.com/naag/terraform-provider-grafanacloud/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestApiKeyCreate(t *testing.T) {
	mock := mock.NewGrafanaCloud().Start()
	mock.HandleApiKeyCreate("my-org", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(grafana.APIKey{Name: "foo"})
	})
	defer mock.Close()

	c, err := portal.NewClient(mock.Server.URL+"/api", "")
	require.NoError(t, err)

	req := &portal.CreateAPIKey{
		Name: "Foo",
		Role: "Admin",
	}
	resp, err := c.CreateAPIKey(req, "my-org")
	require.NoError(t, err)
	require.Equal(t, "foo", resp.Name)
}
*/
