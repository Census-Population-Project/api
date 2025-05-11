package middleware

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/Census-Population-Project/API/internal/service/api/tools"
	"github.com/Census-Population-Project/API/internal/service/auth"

	serviceerrors "github.com/Census-Population-Project/API/internal/errors"

	"github.com/golang-jwt/jwt/v5"
)

func AuthorizationContextSetterMiddleware(authService *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "auth_service", authService))

			next.ServeHTTP(w, r)
		})
	}
}

func AuthorizationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authService, ok := r.Context().Value("auth_service").(*auth.Service)

			if !ok || authService == nil {
				tools.RespondWithError(w, http.StatusInternalServerError, "Service error, sorry")
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				err := NewAuthorizationHeaderIsMissingError()
				tools.RespondWithError(w, err.ErrorStatusCode(), err.Error())
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			access, claims, err := authService.ValidateUserToken(tokenStr)
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

			if !access {
				tools.RespondWithError(w, err.(serviceerrors.ServiceError).ErrorStatusCode(), err.Error())
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "claims", claims))

			next.ServeHTTP(w, r)
		})
	}
}

func RolesMiddleware(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("claims").(*jwt.MapClaims)
			if !ok {
				tools.RespondWithError(w, http.StatusInternalServerError, "Service error, sorry")
				return
			}

			userRole, ok := (*claims)["role"].(string)
			err := serviceerrors.NewForbiddenError()
			if !ok {
				tools.RespondWithError(w, err.ErrorStatusCode(), err.Error())
				return
			}

			if slices.Contains(roles, userRole) {
				next.ServeHTTP(w, r)
			} else {
				tools.RespondWithError(w, err.ErrorStatusCode(), err.Error())
				return
			}
		})
	}
}
