package middleware

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

func LoggingMiddleware(logger *zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info().Msg(fmt.Sprintf("Method: %s, Path: %s", r.Method, r.RequestURI))
		next.ServeHTTP(w, r)
	})
}
