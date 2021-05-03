package mock

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/grafana"
)

func (g *GrafanaCloud) listGrafanaAPIKeys(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "stack")
	sendResponse(w, g.organisation.stackAPIKeys[stackName].Keys, http.StatusOK)
}

func (g *GrafanaCloud) deleteGrafanaAPIKey(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "stack")
	keyID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, err)
		return
	}

	newItems := make([]*grafana.APIKey, 0)
	for _, k := range g.organisation.stackAPIKeys[stackName].Keys {
		if k.ID != keyID {
			newItems = append(newItems, k)
		}
	}

	g.organisation.stackAPIKeys[stackName].Keys = newItems
	sendResponse(w, nil, http.StatusNoContent)
}
