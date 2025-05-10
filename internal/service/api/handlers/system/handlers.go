package system

import (
	"net/http"
	"time"

	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/service/api/response"
	"github.com/Census-Population-Project/API/internal/service/api/tools"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Router *chi.Mux
	Config *config.Config
}

func NewSystemHandler(cfg *config.Config) *Handlers {
	handlers := &Handlers{
		Router: chi.NewRouter(),
		Config: cfg,
	}

	handlers.Router.Get("/ping", handlers.GetPingHandler())
	handlers.Router.Get("/status", handlers.GetStatusHandler())

	return handlers
}

func (h *Handlers) GetPingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tools.RespondWithJSON(w, http.StatusOK, response.SuccessResponse{
			Status: "success",
		})
	}
}

func (h *Handlers) GetStatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tools.RespondWithJSON(w, http.StatusOK, response.SuccessResponseWithResult{
			Status: "success",
			Result: StatusResponse{
				Version:    h.Config.Version,
				ServerTime: time.Now().Unix(),
				DevMode:    h.Config.DevMode,
			},
		})
	}
}
