package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

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
	nextID        int
}

type organisation struct {
	StackList     *portal.StackList
	PortalAPIKeys *portal.APIKeyList
	StackAPIKeys  map[string]*grafana.APIKeyList
}

type errorResponse struct {
	Message string `json:"message"`
}

func (g *GrafanaCloud) Start() *GrafanaCloud {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/api/orgs/{org}", func(r chi.Router) {
		r.Use(g.organisationCtx)
		r.Post("/api-keys", g.createPortalAPIKey)
		r.Get("/api-keys", g.listPortalAPIKeys)
		r.Delete("/api-keys/{name}", g.deletePortalAPIKey)
		r.Get("/instances", g.listStacks)
	})
	r.Post("/api/instances", g.createStack)
	r.Delete("/api/instances/{stackSlug}", g.deleteStack)
	r.Post("/api/instances/{stack}/api/auth/keys", g.createProxyGrafanaAPIKey)

	// Grafana Cloud API doesn't really offer routes at /grafana. These are just provided
	// here so that we can mock the Grafana API running inside Grafana Cloud stacks.
	r.Get("/grafana/{stack}/api/auth/keys", g.listGrafanaAPIKeys)
	r.Delete("/grafana/{stack}/api/auth/keys/{id}", g.deleteGrafanaAPIKey)
	g.Server = httptest.NewServer(r)
	return g
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

func (g *GrafanaCloud) Close() {
	g.Server.Close()
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

func sendResponse(w http.ResponseWriter, v interface{}, status int) {
	if v != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	}

	w.WriteHeader(status)

	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			panic(err)
		}
	}
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

func (g *GrafanaCloud) GetNextID() int {
	g.nextID += 1
	return g.nextID
}
