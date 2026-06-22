package middleware

import (
	"context"
	"net/http"
	"time"
)

func CheckingTimeoutMiddleware(ti int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(ti)*time.Second)
		defer cancel()

		r = r.WithContext(ctx)

		done := make(chan struct{})

		go func() {
			next.ServeHTTP(w, r)
			close(done)
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write([]byte("Request Timeout"))
		}
	})
}
