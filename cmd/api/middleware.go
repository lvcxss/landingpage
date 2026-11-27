package main

import (
	"crypto/subtle"
	"net/http"

	"golang.org/x/time/rate"
)

func rateLimiter(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 10)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded, try again later", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func requireAdminToken(secretToken string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("X-Admin-Token")
		match := subtle.ConstantTimeCompare([]byte(clientToken), []byte(secretToken))
		if match != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
