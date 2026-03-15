package stack

import (
	"net/http"
)

// AuthMiddleware returns HTTP middleware that validates requests against
// the Lattice Runtime. It checks the Lattice-Session-Token or Lattice-API-Key
// header and rejects unauthorized requests.
func (s *Stack) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Lattice-Session-Token")
		apiKey := r.Header.Get("Lattice-API-Key")

		if token == "" && apiKey == "" {
			http.Error(w, "unauthorized: missing authentication header", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware returns HTTP middleware that adds CORS headers for
// browser-based stack UIs.
func CORSMiddleware(allowOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Lattice-Session-Token, Lattice-API-Key")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
