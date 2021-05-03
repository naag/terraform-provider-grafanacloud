package mock

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/grafana"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/portal"
)

func (g *GrafanaCloud) createPortalAPIKey(w http.ResponseWriter, r *http.Request) {
	org := r.Context().Value("organisation").(*organisation)
	apiKey := &portal.APIKey{
		ID:    g.GetNextID(),
		Token: "very-secret",
	}
	fromJSON(apiKey, r)

	org.PortalAPIKeys.Items = append(org.PortalAPIKeys.Items, apiKey)
	sendResponse(w, apiKey, http.StatusCreated)
}

func (g *GrafanaCloud) listPortalAPIKeys(w http.ResponseWriter, r *http.Request) {
	org := r.Context().Value("organisation").(*organisation)
	sendResponse(w, org.PortalAPIKeys, http.StatusOK)
}

func (g *GrafanaCloud) deletePortalAPIKey(w http.ResponseWriter, r *http.Request) {
	org := r.Context().Value("organisation").(*organisation)
	keyName := chi.URLParam(r, "name")

	newItems := make([]*portal.APIKey, 0)
	for _, k := range org.PortalAPIKeys.Items {
		if k.Name != keyName {
			newItems = append(newItems, k)
		}
	}

	org.PortalAPIKeys.Items = newItems
	sendResponse(w, nil, http.StatusNoContent)
}

func (g *GrafanaCloud) listStacks(w http.ResponseWriter, r *http.Request) {
	org := r.Context().Value("organisation").(*organisation)
	sendResponse(w, org.StackList, http.StatusOK)
}

func (g *GrafanaCloud) createStack(w http.ResponseWriter, r *http.Request) {
	orgName := os.Getenv(EnvOrganisation)
	org := g.Organisations[orgName]
	stack := &portal.Stack{
		HmInstancePromID:  g.GetNextID(),
		HmInstancePromURL: "https://prometheus-instance",
		AmInstanceID:      g.GetNextID(),
	}
	fromJSON(stack, r)

	stack.ID = g.GetNextID()
	stack.OrgID = g.GetNextID()
	stack.OrgSlug = orgName
	stack.OrgName = orgName
	if stack.URL == "" {
		stack.URL = fmt.Sprintf("%s/grafana/%s", g.Server.URL, stack.Slug)
	}

	org.StackList.Items = append(org.StackList.Items, stack)
	org.StackAPIKeys = map[string]*grafana.APIKeyList{
		stack.Slug: {},
	}

	sendResponse(w, stack, http.StatusCreated)
}

func (g *GrafanaCloud) deleteStack(w http.ResponseWriter, r *http.Request) {
	orgName := os.Getenv(EnvOrganisation)
	org := g.Organisations[orgName]
	stackSlug := chi.URLParam(r, "stackSlug")

	newItems := make([]*portal.Stack, 0)
	for _, s := range org.StackList.Items {
		if s.Slug != stackSlug {
			newItems = append(newItems, s)
		}
	}

	org.StackList.Items = newItems
	delete(org.StackAPIKeys, stackSlug)
	sendResponse(w, nil, http.StatusNoContent)
}

func (g *GrafanaCloud) createProxyGrafanaAPIKey(w http.ResponseWriter, r *http.Request) {
	orgName := os.Getenv(EnvOrganisation)
	org := g.Organisations[orgName]
	stackName := chi.URLParam(r, "stack")

	apiKey := &grafana.APIKey{
		ID:  g.GetNextID(),
		Key: "very-secret",
	}
	fromJSON(apiKey, r)

	stackAPIKeys := org.StackAPIKeys[stackName]
	stackAPIKeys.Keys = append(stackAPIKeys.Keys, apiKey)
	sendResponse(w, apiKey, http.StatusCreated)
}
