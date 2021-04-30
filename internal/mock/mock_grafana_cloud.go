package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/grafana"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/portal"
)

// TODO: remove me
const EnvOrganisation = "GRAFANA_CLOUD_ORGANISATION"

type GrafanaCloud struct {
	Organisations map[string]*organisation
	Server        *httptest.Server
}

type organisation struct {
	StackList     *portal.StackList
	PortalAPIKeys *portal.APIKeyList
	StackAPIKeys  map[string]*grafana.APIKeyList
	nextID        int
}

type errorResponse struct {
	Message string `json:"message"`
}

func NewGrafanaCloud() *GrafanaCloud {
	return &GrafanaCloud{
		Organisations: make(map[string]*organisation),
	}
}

func (g *GrafanaCloud) WithOrganisation(orgName string) *GrafanaCloud {
	g.Organisations[orgName] = &organisation{
		StackList:     &portal.StackList{},
		PortalAPIKeys: &portal.APIKeyList{},
	}

	return g
}

func (g *GrafanaCloud) WithStack(stackName, orgName string) *GrafanaCloud {
	org := g.Organisations[orgName]
	url := fmt.Sprintf("%s/grafana/%s", g.Server.URL, stackName)

	stack := &portal.Stack{
		ID:      1,
		OrgID:   1,
		OrgSlug: orgName,
		OrgName: orgName,
		Name:    stackName,
		Slug:    stackName,
		URL:     url,
	}

	org.StackList.Items = append(org.StackList.Items, stack)
	org.StackAPIKeys = map[string]*grafana.APIKeyList{
		stackName: {},
	}
	return g
}

func (g *GrafanaCloud) Start() *GrafanaCloud {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/api/orgs/{org}", func(r chi.Router) {
		r.Use(g.organisationCtx)
		r.Post("/api-keys", g.createPortalAPIKey)
		r.Get("/api-keys", g.listPortalAPIKeys)
		r.Delete("/api-keys/{name}", g.deletePortalAPIKey)
		r.Get("/instances", g.listStacks)
	})
	r.Post("/api/instances", g.createStack)
	r.Post("/api/instances/{stack}/api/auth/keys", g.createProxyGrafanaAPIKey)
	r.Get("/grafana/{stack}/api/auth/keys", g.listGrafanaAPIKeys)
	r.Delete("/grafana/{stack}/api/auth/keys/{id}", g.deleteGrafanaAPIKey)
	g.Server = httptest.NewServer(r)
	return g
}

func (g *GrafanaCloud) Close() {
	g.Server.Close()
}

func (g *GrafanaCloud) createPortalAPIKey(w http.ResponseWriter, r *http.Request) {
	org := r.Context().Value("organisation").(*organisation)
	apiKey := &portal.APIKey{
		ID:    org.GetNextID(),
		Token: "very-secret",
	}

	if err := json.NewDecoder(r.Body).Decode(apiKey); err != nil {
		sendError(w, err)
		return
	}

	org.PortalAPIKeys.Items = append(org.PortalAPIKeys.Items, apiKey)
	sendCreated(apiKey, w)
}

func (g *GrafanaCloud) listPortalAPIKeys(w http.ResponseWriter, r *http.Request) {
	org := r.Context().Value("organisation").(*organisation)
	sendSuccess(org.PortalAPIKeys, w)
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
	sendDeleted(w)
}

func (g *GrafanaCloud) listStacks(w http.ResponseWriter, r *http.Request) {
	org := r.Context().Value("organisation").(*organisation)
	sendSuccess(org.StackList, w)
}

func (g *GrafanaCloud) createStack(w http.ResponseWriter, r *http.Request) {
	org := r.Context().Value("organisation").(*organisation)
	stack := &portal.Stack{
		ID: org.GetNextID(),
	}

	fromJSON(stack, r)
	org.StackList.Items = append(org.StackList.Items, stack)
	sendCreated(stack, w)
}

func (g *GrafanaCloud) createProxyGrafanaAPIKey(w http.ResponseWriter, r *http.Request) {
	orgName := os.Getenv(EnvOrganisation)
	org := g.Organisations[orgName]
	stackName := chi.URLParam(r, "stack")

	apiKey := &grafana.APIKey{
		ID:  org.GetNextID(),
		Key: "very-secret",
	}

	if err := json.NewDecoder(r.Body).Decode(apiKey); err != nil {
		sendError(w, err)
		return
	}

	stackAPIKeys := org.StackAPIKeys[stackName]
	stackAPIKeys.Keys = append(stackAPIKeys.Keys, apiKey)
	sendCreated(apiKey, w)
}

func (g *GrafanaCloud) listGrafanaAPIKeys(w http.ResponseWriter, r *http.Request) {
	orgName := os.Getenv(EnvOrganisation)
	org := g.Organisations[orgName]
	stackName := chi.URLParam(r, "stack")
	sendSuccess(org.StackAPIKeys[stackName].Keys, w)
}

func (g *GrafanaCloud) deleteGrafanaAPIKey(w http.ResponseWriter, r *http.Request) {
	orgName := os.Getenv(EnvOrganisation)
	org := g.Organisations[orgName]
	stackName := chi.URLParam(r, "stack")
	keyID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, err)
		return
	}

	newItems := make([]*grafana.APIKey, 0)
	for _, k := range org.StackAPIKeys[stackName].Keys {
		if k.ID != keyID {
			newItems = append(newItems, k)
		}
	}

	org.StackAPIKeys[stackName].Keys = newItems
	sendDeleted(w)
}

func fromJSON(d interface{}, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, d); err != nil {
		panic(err)
	}
}

func (g *GrafanaCloud) organisationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orgName := chi.URLParam(r, "org")

		org, err := g.findOrganisation(orgName)
		if err != nil {
			sendError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "organisation", org)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (g *GrafanaCloud) findOrganisation(org string) (*organisation, error) {
	o, ok := g.Organisations[org]
	if !ok {
		return nil, fmt.Errorf("failed to find organisation `%s`", org)
	}

	return o, nil
}

func sendSuccess(d interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		panic(err)
	}
}

func sendCreated(d interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(d); err != nil {
		panic(err)
	}
}

func sendDeleted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func sendError(w http.ResponseWriter, err error) {
	resp := &errorResponse{
		Message: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func (o *organisation) GetNextID() int {
	o.nextID += 1
	return o.nextID
}
