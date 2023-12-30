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
	client, err := ent.Open("sqlite3", "file:bib.db?cache=shared&_fk=1")
	if err != nil {
		panic(fmt.Sprintf("Can't open database: %s\n", err))
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		panic(fmt.Sprintf("Can't create schema: %s\n", err))
	}

	defer client.Close()

	app := &App{client, log.Default()}
	r := mux.NewRouter()
	r.Use(func(h http.Handler) http.Handler { return handlers.LoggingHandler(os.Stdout, h) })
	//r.Use(handlers.RecoveryHandler())

	authnz, err := NewAuth(context.Background(), "http://127.0.0.1:5556/dex", "example-app", "ZXhhbXBsZS1hcHAtc2VjcmV0", "http://localhost:5555/auth/callback", "http://localhost:8080")
	if err != nil {
		panic(fmt.Sprintf("Can't configure auth: %s", err.Error()))
	}
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", authnz.LoginHandler)
	auth.HandleFunc("/callback", authnz.CallbackHandler)

	api := r.PathPrefix("/api/v1").Subrouter()
	corsenv := os.Getenv("DEV_CORS")
	if corsenv != "" {
		app.log.Printf("Using %s as CORS origin host\n", corsenv)
		cors := &CORSMiddleWare{
			AllowedOrigins: []string{corsenv},
			AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		}
		api.Use(cors.Handler)
	}

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

	var csrfkey []byte
	csrfenv := os.Getenv("CSRF_KEY")
	if csrfenv == "" {
		csrfkey = []byte("really-should-change-this")
	} else {
		csrfkey = []byte(csrfenv)
	}
	CSRF := csrf.Protect(csrfkey)
	api.Use(CSRF)

	r.PathPrefix("/app/").Handler(http.StripPrefix("/app/", http.FileServer(http.Dir("client"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:5555",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
