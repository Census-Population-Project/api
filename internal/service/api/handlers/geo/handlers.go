package geo

import (
	"errors"
	"math"
	"net/http"

	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/service/api/middleware"
	"github.com/Census-Population-Project/API/internal/service/api/response"
	"github.com/Census-Population-Project/API/internal/service/api/tools"
	"github.com/Census-Population-Project/API/internal/service/geo"

	serviceerrors "github.com/Census-Population-Project/API/internal/errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handlers struct {
	Router *chi.Mux
	Config *config.Config

	GeoService *geo.Service
}

func NewGeoHandler(cfg *config.Config, geoService *geo.Service) *Handlers {
	handlers := &Handlers{
		Router: chi.NewRouter(),
		Config: cfg,

		GeoService: geoService,
	}

	handlers.Router.With(middleware.AuthorizationMiddleware()).Post("/suggestions", handlers.PostAddressSuggestionsHandler())

	handlers.Router.With(middleware.AuthorizationMiddleware()).Get("/regions", handlers.GetRegionsHandler())
	handlers.Router.With(middleware.AuthorizationMiddleware()).Get("/regions/{region-id}/cities", handlers.GetCitiesInRegionHandler())

	handlers.Router.With(middleware.AuthorizationMiddleware()).Get("/cities", handlers.GetCitiesHandler())
	handlers.Router.With(middleware.AuthorizationMiddleware()).Get("/cities/{city-id}/buildings", handlers.GetBuildingsInCityHandler())

	handlers.Router.With(middleware.AuthorizationMiddleware()).Get("/buildings", handlers.GetBuildingsHandler())
	handlers.Router.With(middleware.AuthorizationMiddleware()).Get("/buildings/{building-id}/addresses", handlers.GetAddressesInBuildingHandler())

	return handlers
}

func (h *Handlers) PostAddressSuggestionsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			tools.RespondWithError(w, http.StatusBadRequest, "Address is required")
			return
		}

		limit, err := tools.ParseIntQuery(r, "limit", 1, 20, 10, false)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		suggestionsValuesData, err := h.GeoService.AddressSuggestionsValues(address, *limit)
		if err != nil {
			var srvErr serviceerrors.ServiceError
			if errors.As(err, &srvErr) {
				tools.RespondWithError(w, srvErr.ErrorStatusCode(), err.Error())
			} else {
				tools.RespondWithError(w, http.StatusInternalServerError, "Service error, sorry")
			}
			return
		}

		tools.RespondWithJSON(w, http.StatusOK, response.SuccessResponseWithResult{
			Status: "success",
			Result: suggestionsValuesData,
		})
	}
}

func (h *Handlers) GetRegionsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := tools.ParseIntQuery(r, "limit", 0, 10, 10, false)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		offset, err := tools.ParseIntQuery(r, "offset", 0, math.MaxInt, 0, false)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		regionsData, err := h.GeoService.GetRegions(*limit, *offset)
		if err != nil {
			var srvErr serviceerrors.ServiceError
			if errors.As(err, &srvErr) {
				tools.RespondWithError(w, srvErr.ErrorStatusCode(), err.Error())
			} else {
				tools.RespondWithError(w, http.StatusInternalServerError, "Service error, sorry")
			}
			return
		}

		tools.RespondWithJSON(w, http.StatusOK, response.SuccessResponseWithResult{
			Status: "success",
			Result: regionsData,
		})
	}
}

func (h *Handlers) GetCitiesInRegionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		regionIdStr := chi.URLParam(r, "region-id")
		if regionIdStr == "" {
			tools.RespondWithError(w, http.StatusBadRequest, "Region id is required")
			return
		}

		regionId, err := uuid.Parse(regionIdStr)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, "Invalid region id")
			return
		}

		limit, err := tools.ParseIntQuery(r, "limit", 0, 10, 10, false)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		offset, err := tools.ParseIntQuery(r, "offset", 0, math.MaxInt, 0, false)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		regionData, err := h.GeoService.GetCitiesInRegion(regionId, *limit, *offset)
		if err != nil {
			var srvErr serviceerrors.ServiceError
			if errors.As(err, &srvErr) {
				tools.RespondWithError(w, srvErr.ErrorStatusCode(), err.Error())
			} else {
				tools.RespondWithError(w, http.StatusInternalServerError, "Service error, sorry")
			}
			return
		}

		tools.RespondWithJSON(w, http.StatusOK, response.SuccessResponseWithResult{
			Status: "success",
			Result: regionData,
		})
	}
}

func (h *Handlers) GetCitiesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := tools.ParseIntQuery(r, "limit", 0, 10, 10, false)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		offset, err := tools.ParseIntQuery(r, "offset", 0, math.MaxInt, 0, false)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		citiesData, err := h.GeoService.GetCities(*limit, *offset)
		if err != nil {
			var srvErr serviceerrors.ServiceError
			if errors.As(err, &srvErr) {
				tools.RespondWithError(w, srvErr.ErrorStatusCode(), err.Error())
			} else {
				tools.RespondWithError(w, http.StatusInternalServerError, "Service error, sorry")
			}
			return
		}

		tools.RespondWithJSON(w, http.StatusOK, response.SuccessResponseWithResult{
			Status: "success",
			Result: citiesData,
		})
	}
}

func (h *Handlers) GetBuildingsInCityHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cityIdStr := chi.URLParam(r, "city-id")
		if cityIdStr == "" {
			tools.RespondWithError(w, http.StatusBadRequest, "City id is required")
			return
		}

		cityId, err := uuid.Parse(cityIdStr)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, "Invalid city id")
			return
		}

		limit, err := tools.ParseIntQuery(r, "limit", 0, 10, 10, false)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		offset, err := tools.ParseIntQuery(r, "offset", 0, math.MaxInt, 0, false)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		cityData, err := h.GeoService.GetBuildingsInCity(cityId, *limit, *offset)
		if err != nil {
			var srvErr serviceerrors.ServiceError
			if errors.As(err, &srvErr) {
				tools.RespondWithError(w, srvErr.ErrorStatusCode(), err.Error())
			} else {
				tools.RespondWithError(w, http.StatusInternalServerError, "Service error, sorry")
			}
			return
		}

		tools.RespondWithJSON(w, http.StatusOK, response.SuccessResponseWithResult{
			Status: "success",
			Result: cityData,
		})
	}
}

func (h *Handlers) GetBuildingsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := tools.ParseIntQuery(r, "limit", 0, 10, 10, false)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		offset, err := tools.ParseIntQuery(r, "offset", 0, math.MaxInt, 0, false)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		buildingsData, err := h.GeoService.GetBuildings(*limit, *offset)
		if err != nil {
			var srvErr serviceerrors.ServiceError
			if errors.As(err, &srvErr) {
				tools.RespondWithError(w, srvErr.ErrorStatusCode(), err.Error())
			} else {
				tools.RespondWithError(w, http.StatusInternalServerError, "Service error, sorry")
			}
			return
		}

		tools.RespondWithJSON(w, http.StatusOK, response.SuccessResponseWithResult{
			Status: "success",
			Result: buildingsData,
		})
	}
}

func (h *Handlers) GetAddressesInBuildingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buildingIdStr := chi.URLParam(r, "building-id")
		if buildingIdStr == "" {
			tools.RespondWithError(w, http.StatusBadRequest, "Building id is required")
			return
		}

		buildingId, err := uuid.Parse(buildingIdStr)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, "Invalid building id")
			return
		}

		limit, err := tools.ParseIntQuery(r, "limit", 0, 10, 10, false)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		offset, err := tools.ParseIntQuery(r, "offset", 0, math.MaxInt, 0, false)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		buildingData, err := h.GeoService.GetAddressesInBuilding(buildingId, *limit, *offset)
		if err != nil {
			var srvErr serviceerrors.ServiceError
			if errors.As(err, &srvErr) {
				tools.RespondWithError(w, srvErr.ErrorStatusCode(), err.Error())
			} else {
				tools.RespondWithError(w, http.StatusInternalServerError, "Service error, sorry")
			}
			return
		}

		tools.RespondWithJSON(w, http.StatusOK, response.SuccessResponseWithResult{
			Status: "success",
			Result: buildingData,
		})
	}
}
