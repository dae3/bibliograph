package main

import (
	"fmt"
	"net/http"
	"strings"
)

/*

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)
*/

const (
	ParamBookId = "bookid"
)

func VerifyContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			ct := r.Header["Content-Type"][0] // anyone sending multiple content-type headers gets what they deserve
			if ct == "application/json" {
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
				fmt.Printf("CORS says no")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
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

/*
// bookid_param attempts to parse an integer path parameter named "id".
// If successful the int value of the parameter is added to the context
// with the key "bookid", otherwise http.StatusBadRequest is returned
func GetBookIdParam(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idparam, ok := vars["id"]
	if !ok {
		return 0, new error()
		http.Error(w, "Missing id path parameter", http.StatusBadRequest)
	} else {
		id, err := strconv.ParseInt(idparam, 0, 0)
		if err != nil {
			http.Error(w, "Can't parse id path parameter", http.StatusBadRequest)
		} else {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ParamBookId, id)))
		}
	}
}
*/
