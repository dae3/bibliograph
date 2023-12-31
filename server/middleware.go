package main

import (
	"net/http"
	"strings"
)

func VerifyContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			ct, ok := r.Header["Content-Type"]
			if ok && ct[0] == "application/json" { // anyone sending multiple content-type headers gets what they deserve
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "", http.StatusUnsupportedMediaType)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type CORSMiddleWare struct {
	AllowedOrigins []string
	AllowedMethods []string
}

func (m *CORSMiddleWare) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "" { // not CORS
			next.ServeHTTP(w, r)
		} else {
			for _, v := range m.AllowedOrigins {
				if v == origin {
					break
				}
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join([]string{"Content-Type", "X-CSRF-Token"}, ", "))
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			} else {
				next.ServeHTTP(w, r)
			}
		}
	})
}
