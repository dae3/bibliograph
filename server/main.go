package main

import (
	"bibliograph/api/ent"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	db  *ent.Client
	log *log.Logger
}

func main() {
	conf, err := ParseConfig()
	if err != nil {
		panic(err.Error())
	}

	client, err := ent.Open(conf.DatabaseDriver, conf.DatabaseConnectionString)
	if err != nil {
		panic(fmt.Sprintf("Can't open database: %s\n", err))
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		panic(fmt.Sprintf("Can't create schema: %s\n", err))
	}
	defer client.Close()

	var pap string
	if conf.CORSOrigin != "" {
		pap = conf.CORSOrigin
	} else {
		pap = "/"
	}
	authnz, err := NewAuth(context.Background(), &conf.OIDCDiscovery, conf.OIDCClientID, conf.OIDCClientSecret, &conf.OIDCRedirectURL, pap, conf.SessionStoreKey)
	if err != nil {
		panic(fmt.Sprintf("Can't configure auth: %s", err.Error()))
	}

	app := &App{client, log.Default()}
	r := mux.NewRouter()
	r.Use(func(h http.Handler) http.Handler { return handlers.LoggingHandler(os.Stdout, h) })
	r.Use(handlers.RecoveryHandler())

	// auth routes
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", authnz.LoginHandler)
	auth.HandleFunc("/logout", authnz.LogoutHandler)
	auth.HandleFunc("/callback", authnz.CallbackHandler)
	auth.HandleFunc("/status", authnz.StatusHandler)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(VerifyContentTypeMiddleware)
	api.HandleFunc("/books", app.books_get).Methods(http.MethodGet)
	api.HandleFunc("/books/{id}", app.book_get).Methods(http.MethodGet)
	api.Handle("/books", authnz.AuthorizerMiddleware(http.HandlerFunc(app.book_post_new))).Methods(http.MethodPost, http.MethodOptions)
	api.Handle("/books/{id}", authnz.AuthorizerMiddleware(http.HandlerFunc(app.book_post))).Methods(http.MethodPost, http.MethodOptions)
	api.Handle("/books/{id}/refs", authnz.AuthorizerMiddleware(http.HandlerFunc(app.ref_post))).Methods(http.MethodPost, http.MethodOptions)
	api.Handle("/books/{id}", authnz.AuthorizerMiddleware(http.HandlerFunc(app.book_delete))).Methods(http.MethodDelete, http.MethodOptions)
	api.Handle("/books/{id}/refs/{refid}", authnz.AuthorizerMiddleware(http.HandlerFunc(app.ref_delete))).Methods(http.MethodDelete, http.MethodOptions)
	api.Handle("/token", authnz.AuthorizerMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csrf.Token(r)))
	}))).Methods(http.MethodGet, http.MethodOptions) // CORS preflight because of unsafe header

	if conf.CORSOrigin != "" {
		app.log.Printf("Using %s as CORS origin host\n", conf.CORSOrigin)
		cors := &CORSMiddleWare{
			AllowedOrigins: []string{conf.CORSOrigin},
			AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		}
		api.Use(cors.Handler)
		auth.Use(cors.Handler)
	}

	CSRF := csrf.Protect([]byte(conf.CSRFKey))
	api.Use(CSRF)

	// static hosting of SPA
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusPermanentRedirect)
	})
	r.PathPrefix("/app/").Handler(http.StripPrefix("/app/", http.FileServer(http.FS(SpaFS))))

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", conf.ListenPort),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
