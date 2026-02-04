package main

import "net/http"

// Renamed from cacheMiddleware to noCacheMiddleware
func noCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Updated to set Cache-Control header to no-store
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}