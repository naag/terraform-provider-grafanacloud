package mock

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/grafana"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/portal"
)

func (g *GrafanaCloud) createPortalAPIKey(w http.ResponseWriter, r *http.Request) {
	apiKey := &portal.APIKey{
		ID:    g.GetNextID(),
		Token: "very-secret",
	}
	fromJSON(apiKey, r)

	g.organisation.portalAPIKeys.Items = append(g.organisation.portalAPIKeys.Items, apiKey)
	sendResponse(w, apiKey, http.StatusCreated)
}

func (g *GrafanaCloud) listPortalAPIKeys(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, g.organisation.portalAPIKeys, http.StatusOK)
}

func (g *GrafanaCloud) deletePortalAPIKey(w http.ResponseWriter, r *http.Request) {
	keyName := chi.URLParam(r, "name")

	newItems := make([]*portal.APIKey, 0)
	for _, k := range g.organisation.portalAPIKeys.Items {
		if k.Name != keyName {
			newItems = append(newItems, k)
		}
	}

	g.organisation.portalAPIKeys.Items = newItems
	sendResponse(w, nil, http.StatusNoContent)
}

func (g *GrafanaCloud) listStacks(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, g.organisation.stackList, http.StatusOK)
}

func (g *GrafanaCloud) createStack(w http.ResponseWriter, r *http.Request) {
	stack := &portal.Stack{
		HmInstancePromID:  g.GetNextID(),
		HmInstancePromURL: "https://prometheus-instance",
		AmInstanceID:      g.GetNextID(),
	}
	fromJSON(stack, r)

	stack.ID = g.GetNextID()
	stack.OrgID = g.GetNextID()
	stack.OrgSlug = g.organisation.name
	stack.OrgName = g.organisation.name
	if stack.URL == "" {
		stack.URL = fmt.Sprintf("%s/grafana/%s", g.URL(), stack.Slug)
	}

	g.organisation.stackList.Items = append(g.organisation.stackList.Items, stack)
	g.organisation.stackAPIKeys[stack.Slug] = &grafana.APIKeyList{}

	sendResponse(w, stack, http.StatusCreated)
}

func (g *GrafanaCloud) deleteStack(w http.ResponseWriter, r *http.Request) {
	stackSlug := chi.URLParam(r, "stack")

	newItems := make([]*portal.Stack, 0)
	for _, s := range g.organisation.stackList.Items {
		if s.Slug != stackSlug {
			newItems = append(newItems, s)
		}
	}

	g.organisation.stackList.Items = newItems
	delete(g.organisation.stackAPIKeys, stackSlug)
	sendResponse(w, nil, http.StatusNoContent)
}

func (g *GrafanaCloud) createGrafanaAPIKeyProxy(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "stack")

	apiKey := &grafana.APIKey{
		ID:  g.GetNextID(),
		Key: "very-secret",
	}
	fromJSON(apiKey, r)

	stackAPIKeys := g.organisation.stackAPIKeys[stackName]
	stackAPIKeys.Keys = append(stackAPIKeys.Keys, apiKey)
	sendResponse(w, apiKey, http.StatusCreated)
}
