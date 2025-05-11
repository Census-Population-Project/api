package users

import (
	"errors"
	"math"
	"net/http"

	"github.com/Census-Population-Project/API/internal/config"
	"github.com/Census-Population-Project/API/internal/service/api/middleware"
	"github.com/Census-Population-Project/API/internal/service/api/response"
	"github.com/Census-Population-Project/API/internal/service/api/tools"
	"github.com/Census-Population-Project/API/internal/service/users"

	serviceerrors "github.com/Census-Population-Project/API/internal/errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handlers struct {
	Router *chi.Mux
	Config *config.Config

	UsersService *users.Service
}

func NewUsersHandler(cfg *config.Config, usersService *users.Service) *Handlers {
	handlers := &Handlers{
		Router: chi.NewRouter(),
		Config: cfg,

		UsersService: usersService,
	}

	handlers.Router.With(middleware.AuthorizationMiddleware(), middleware.RolesMiddleware("administrator")).Get("/", handlers.GetUsersHandler())
	handlers.Router.With(middleware.AuthorizationMiddleware()).Get("/me", handlers.GetUserMeHandler())
	handlers.Router.With(middleware.AuthorizationMiddleware(), middleware.RolesMiddleware("administrator")).Get("/{user-id}", handlers.GetUserHandler())
	handlers.Router.With(middleware.AuthorizationMiddleware(), middleware.RolesMiddleware("administrator")).Post("/", handlers.PostCreateUserHandler())
	handlers.Router.With(middleware.AuthorizationMiddleware(), middleware.RolesMiddleware("administrator")).Patch("/{user-id}", handlers.PatchUserHandler())

	return handlers
}

func (h *Handlers) GetUsersHandler() http.HandlerFunc {
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

		usersData, err := h.UsersService.GetUsers(*limit, *offset)
		if err != nil {
			var srvErr serviceerrors.ServiceError
			if errors.As(err, &srvErr) {
				tools.RespondWithError(w, srvErr.ErrorStatusCode(), srvErr.Error())
				return
			} else {
				tools.RespondWithError(w, http.StatusInternalServerError, "Service error, sorry")
				return
			}
		}

		tools.RespondWithJSON(w, http.StatusOK, response.SuccessResponseWithResult{
			Status: "success",
			Result: usersData,
		})
	}
}

func (h *Handlers) GetUserMeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := tools.GetUserIdFromContext(r)
		if err != nil {
			tools.RespondWithError(w, http.StatusInternalServerError, "Service error, sorry")
			return
		}

		userData, err := h.UsersService.GetUserByID(*userId)
		if err != nil {
			var srvErr serviceerrors.ServiceError
			if errors.As(err, &srvErr) {
				tools.RespondWithError(w, srvErr.ErrorStatusCode(), srvErr.Error())
				return
			} else {
				tools.RespondWithError(w, http.StatusInternalServerError, "Service error, sorry")
				return
			}
		}

		tools.RespondWithJSON(w, http.StatusOK, response.SuccessResponseWithResult{
			Status: "success",
			Result: userData,
		})
	}
}

func (h *Handlers) GetUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdStr := chi.URLParam(r, "user-id")
		if userIdStr == "" {
			tools.RespondWithError(w, http.StatusNotFound, "User id is required")
			return
		}

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, "Invalid user id")
			return
		}

		userData, err := h.UsersService.GetUserByID(userId)
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
			Result: userData,
		})
	}
}

func (h *Handlers) PostCreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := tools.DecodeJSON(w, r, &req); err != nil {
			var srvErr serviceerrors.ServiceError
			if errors.As(err, &srvErr) {
				tools.RespondWithError(w, srvErr.ErrorStatusCode(), srvErr.Error())
			} else {
				tools.RespondWithError(w, http.StatusBadRequest, "Invalid request")
			}
			return
		}

		userData, err := h.UsersService.CreateUser(req.Email, req.Password, req.FirstName, req.LastName, req.Role, false)
		if err != nil {
			var srvErr serviceerrors.ServiceError
			if errors.As(err, &srvErr) {
				tools.RespondWithError(w, srvErr.ErrorStatusCode(), err.Error())
			} else {
				tools.RespondWithError(w, http.StatusInternalServerError, "Service error, sorry")
			}
			return
		}

		tools.RespondWithJSON(w, http.StatusCreated, response.SuccessResponseWithResult{
			Status: "success",
			Result: userData.ID,
		})
	}
}

func (h *Handlers) PatchUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdStr := chi.URLParam(r, "user-id")
		if userIdStr == "" {
			tools.RespondWithError(w, http.StatusNotFound, "User id is required")
			return
		}

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			tools.RespondWithError(w, http.StatusBadRequest, "Invalid user id")
			return
		}

		var req UpdateUserRequest
		if err := tools.DecodeJSON(w, r, &req); err != nil {
			var srvErr serviceerrors.ServiceError
			if errors.As(err, &srvErr) {
				tools.RespondWithError(w, srvErr.ErrorStatusCode(), err.Error())
			} else {
				tools.RespondWithError(w, http.StatusBadRequest, "Invalid request")
			}
			return
		}

		userData, err := h.UsersService.UpdateUserByID(
			userId,
			req.Email,
			req.Password,
			req.FirstName,
			req.LastName,
			req.Role,
			true,
		)
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
			Result: userData,
		})
	}
}
