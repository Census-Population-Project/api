package auth

import (
	"errors"
	"net/http"

	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/service/api/response"
	"github.com/Census-Population-Project/API/internal/service/api/tools"
	"github.com/Census-Population-Project/API/internal/service/auth"

	serviceerrors "github.com/Census-Population-Project/API/internal/errors"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Router *chi.Mux
	Config *config.Config

	AuthService *auth.Service
}

func NewAuthHandler(cfg *config.Config, authService *auth.Service) *Handlers {
	handlers := &Handlers{
		Router: chi.NewRouter(),
		Config: cfg,

		AuthService: authService,
	}

	handlers.Router.Post("/login", handlers.PostLoginHandler())

	return handlers
}

func (h *Handlers) PostLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := tools.DecodeJSON(w, r, &req); err != nil {
			var srvErr serviceerrors.ServiceError
			if errors.As(err, &srvErr) {
				tools.RespondWithError(w, srvErr.ErrorStatusCode(), srvErr.Error())
			} else {
				tools.RespondWithError(w, http.StatusBadRequest, "Invalid request")
			}
			return
		}

		tokensData, err := h.AuthService.LoginUser(req.Email, req.Password)
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
			Result: tokensData,
		})
	}
}
