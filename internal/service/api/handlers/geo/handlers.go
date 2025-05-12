package geo

import (
	"net/http"

	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/service/api/middleware"
	"github.com/Census-Population-Project/API/internal/service/geo"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Router *chi.Mux
	Config *config.Config

	GeoService *geo.Service
}

func NewUsersHandler(cfg *config.Config, geoService *geo.Service) *Handlers {
	handlers := &Handlers{
		Router: chi.NewRouter(),
		Config: cfg,

		GeoService: geoService,
	}

	handlers.Router.With(middleware.AuthorizationMiddleware()).Get("/regions", handlers.GetRegionsHandler())

	return handlers
}

func (handlers *Handlers) GetRegionsHandler() http.HandlerFunc { // TODO: Implement.
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
