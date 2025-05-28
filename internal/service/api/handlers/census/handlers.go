package census

import (
	"errors"
	"math"
	"net/http"

	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/service/api/response"
	"github.com/Census-Population-Project/API/internal/service/api/tools"
	"github.com/Census-Population-Project/API/internal/service/census"
	"github.com/Census-Population-Project/API/internal/service/geo"

	serviceerrors "github.com/Census-Population-Project/API/internal/errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handlers struct {
	Router *chi.Mux
	Config *config.Config

	CensusService *census.Service
	GeoService    *geo.Service
}

func NewCensusHandler(cfg *config.Config, censusService *census.Service, geoService *geo.Service) *Handlers {
	handlers := &Handlers{
		Router: chi.NewRouter(),
		Config: cfg,

		CensusService: censusService,
		GeoService:    geoService,
	}

	handlers.Router.Get("/events", handlers.GetEventsHandler())
	handlers.Router.Get("/events/{event-id}", handlers.GetEventInfoHandler())
	handlers.Router.Get("/events/{event-id}/statistics", handlers.GetEventStatisticsHandler())

	handlers.Router.Get("/events/{event-id}/region/{region-id}", handlers.GetEventDataInRegionHandler())
	handlers.Router.Get("/events/{event-id}/region/{region-id}/statistics", handlers.GetEventStatisticsInRegionHandler())

	handlers.Router.Get("/events/{event-id}/city/{city-id}", handlers.GetEventDataInCityHandler())
	handlers.Router.Get("/events/{event-id}/city/{city-id}/statistics", handlers.GetEventStatisticsInCityHandler())

	return handlers
}

func (h *Handlers) GetEventsHandler() http.HandlerFunc {
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

		regionsData, regionsCount, err := h.CensusService.GetEvents(*limit, *offset)
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
			Result: census.Events{
				Events: regionsData,
				Total:  *regionsCount,
			},
		})
	}
}

func (h *Handlers) GetEventInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventIdStr := chi.URLParam(r, "event-id")
		if eventIdStr == "" {
			tools.RespondWithError(w, http.StatusBadRequest, "Event id is required")
			return
		}

		eventId, err := uuid.Parse(eventIdStr)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, "Invalid event id")
			return
		}

		eventData, err := h.CensusService.GetEventInfoByID(eventId)
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
			Result: eventData,
		})
	}
}

func (h *Handlers) GetEventStatisticsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventIdStr := chi.URLParam(r, "event-id")
		if eventIdStr == "" {
			tools.RespondWithError(w, http.StatusBadRequest, "Event id is required")
			return
		}

		eventId, err := uuid.Parse(eventIdStr)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, "Invalid event id")
			return
		}

		eventStatisticsAll, err := h.CensusService.GetEventStatisticsAllByID(eventId)
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
			Result: eventStatisticsAll,
		})
	}
}

func (h *Handlers) GetEventDataInRegionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventIdStr := chi.URLParam(r, "event-id")
		if eventIdStr == "" {
			tools.RespondWithError(w, http.StatusBadRequest, "Event id is required")
			return
		}

		eventId, err := uuid.Parse(eventIdStr)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, "Invalid event id")
			return
		}

		var regionId *uuid.UUID
		if regionIdStr := chi.URLParam(r, "region-id"); regionIdStr != "" {
			id, err := uuid.Parse(regionIdStr)
			if err != nil {
				tools.RespondWithError(w, http.StatusBadRequest, "Invalid region id")
				return
			}
			regionId = &id
		}

		eventData, err := h.CensusService.GetEventInfoByLocationIDs(eventId, regionId, nil)
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
			Result: eventData,
		})
	}
}

func (h *Handlers) GetEventStatisticsInRegionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventIdStr := chi.URLParam(r, "event-id")
		if eventIdStr == "" {
			tools.RespondWithError(w, http.StatusBadRequest, "Event id is required")
			return
		}

		eventId, err := uuid.Parse(eventIdStr)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, "Invalid event id")
			return
		}

		var regionId *uuid.UUID
		if regionIdStr := chi.URLParam(r, "region-id"); regionIdStr != "" {
			id, err := uuid.Parse(regionIdStr)
			if err != nil {
				tools.RespondWithError(w, http.StatusBadRequest, "Invalid region id")
				return
			}
			regionId = &id
		}

		eventStatistics, err := h.CensusService.GetEventStatisticsByLocationIDs(eventId, regionId, nil)
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
			Result: eventStatistics,
		})
	}
}

func (h *Handlers) GetEventDataInCityHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventIdStr := chi.URLParam(r, "event-id")
		if eventIdStr == "" {
			tools.RespondWithError(w, http.StatusBadRequest, "Event id is required")
			return
		}

		eventId, err := uuid.Parse(eventIdStr)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, "Invalid event id")
			return
		}

		var cityId *uuid.UUID
		if cityIdStr := chi.URLParam(r, "city-id"); cityIdStr != "" {
			id, err := uuid.Parse(cityIdStr)
			if err != nil {
				tools.RespondWithError(w, http.StatusBadRequest, "Invalid city id")
				return
			}
			cityId = &id
		}

		eventData, err := h.CensusService.GetEventInfoByLocationIDs(eventId, nil, cityId)
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
			Result: eventData,
		})
	}
}

func (h *Handlers) GetEventStatisticsInCityHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventIdStr := chi.URLParam(r, "event-id")
		if eventIdStr == "" {
			tools.RespondWithError(w, http.StatusBadRequest, "Event id is required")
			return
		}

		eventId, err := uuid.Parse(eventIdStr)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, "Invalid event id")
			return
		}

		var cityId *uuid.UUID
		if cityIdStr := chi.URLParam(r, "city-id"); cityIdStr != "" {
			id, err := uuid.Parse(cityIdStr)
			if err != nil {
				tools.RespondWithError(w, http.StatusBadRequest, "Invalid city id")
				return
			}
			cityId = &id
		}

		eventStatistics, err := h.CensusService.GetEventStatisticsByLocationIDs(eventId, nil, cityId)
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
			Result: eventStatistics,
		})
	}
}
