package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func HttpLoggerMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			defer func() {
				logger.WithFields(logrus.Fields{
					"method":   r.Method,
					"url":      r.URL.Path,
					"host":     r.Host,
					"remote":   r.RemoteAddr,
					"duration": time.Since(startTime),
				}).Info("Request received")
			}()

			next.ServeHTTP(w, r)
		})
	}
}
