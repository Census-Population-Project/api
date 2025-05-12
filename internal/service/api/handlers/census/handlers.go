package census

import (
	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/service/geo"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Router *chi.Mux
	Config *config.Config

	GeoService *geo.Service
}

func NewCensusHandler(cfg *config.Config, geoService *geo.Service) *Handlers {
	handlers := &Handlers{
		Router: chi.NewRouter(),
		Config: cfg,

		GeoService: geoService,
	}

	return handlers
}
