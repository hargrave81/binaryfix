package main

import (
	"net/http"
	"strings"
)

// Define our struct
type authenticationMiddleware struct {
	tokenUsers map[string]string
}

// Initialize it somewhere
func (amw *authenticationMiddleware) Populate() {

}

// Middleware function, which will be called for each request
func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("session")
		if err != nil || token == nil {
			http.Error(w, "Forbidden", 403)
		}
		if value, err := Decrypt(token.Value); err == nil {
			// We found the token in our map
			user := strings.Split(value, "-")[1]
			user = strings.Split(user, "/")[0]
			r.SetBasicAuth(user, "ok")
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", 403)
		}
	})
}
