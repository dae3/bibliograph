package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

const (
	SessionCookie = "session"
)

type Auth struct {
	RedirectUrl  string
	PostAuthUrl  string
	SessionStore *sessions.CookieStore

	IdP          *oidc.Provider
	OAuth2Config *oauth2.Config
	Verifier     *oidc.IDTokenVerifier
}

func NewAuth(ctx context.Context, IdPUrl string, ClientId string, ClientSecret string, CallbackUrl string, PostAuthUrl string) (auth *Auth, err error) {
	idp, err := oidc.NewProvider(ctx, IdPUrl)
	if err != nil {
		return
	}

	cskey := os.Getenv("SESSION_STORE_KEY")
	if cskey == "" {
		cskey = "s00pers3cr3t"
	}

	auth = &Auth{
		RedirectUrl:  CallbackUrl,
		PostAuthUrl:  PostAuthUrl,
		SessionStore: sessions.NewCookieStore([]byte(cskey)),
		IdP:          idp,
		OAuth2Config: &oauth2.Config{
			Endpoint:     idp.Endpoint(),
			ClientID:     ClientId,
			ClientSecret: ClientSecret,
			RedirectURL:  CallbackUrl,
			Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		},
		Verifier: idp.Verifier(&oidc.Config{ClientID: ClientId}),
	}

	auth.SessionStore.Options.Path = "/"
	auth.SessionStore.Options.Secure = true
	auth.SessionStore.Options.HttpOnly = true

	return
}

func (a *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.OAuth2Config.AuthCodeURL("notveryrandom"), http.StatusTemporaryRedirect)
}

func (a *Auth) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if state != "notveryrandom" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	token, err := a.OAuth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	rawidtoken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	idtoken, err := a.Verifier.Verify(r.Context(), rawidtoken)

	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}

	if err := idtoken.Claims(&claims); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	s, err := a.SessionStore.Get(r, SessionCookie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	s.Values["name"] = claims.Email
	err = s.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, a.PostAuthUrl, http.StatusTemporaryRedirect)
	}
	return
}

func (a *Auth) AuthorizerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := a.SessionStore.Get(r, SessionCookie)
		if err != nil {
			http.Error(w, fmt.Sprintf("Can't get session store: %s", err.Error()), http.StatusInternalServerError)
		}

		_, hasemail := s.Values["name"]

		if hasemail {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "", http.StatusForbidden)
		}
	})
}
