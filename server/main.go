package main

import (
	"bibliograph/api/ent"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	oidcauth "github.com/TJM/gin-gonic-oidcauth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"

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
	g := gin.Default()

	if corsenv := os.Getenv("DEV_CORS"); corsenv != "" {
		app.log.Printf("Using %s as CORS origin host\n", corsenv)
		g.Use(cors.New(cors.Options{
			AllowedOrigins: []string{corsenv},
			AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		}))
	}

	store := cookie.NewStore([]byte("secret"))
	g.Use(sessions.Sessions("mysession", store))
	auth, err := oidcauth.GetOidcAuth(oidcauth.ExampleConfigDex())
	if err != nil {
		panic(fmt.Sprintf("Auth setup failed: %s\n", err))
	}
	g.GET("/login", auth.Login)
	g.GET("/callback", auth.AuthCallback)
	g.GET("/logout", auth.Logout)

	api := g.Group("/api/v1")
	api.GET("/books", app.books_get)
	api.GET("/books/:id", bookid_param, app.book_get)
	// api.POST("/books", auth.AuthRequired(), app.book_post)
	api.POST("/books", app.book_post)
	api.POST("/books/:id/refs", bookid_param, app.ref_post)
	api.DELETE("/books/:id", bookid_param, app.book_delete)
	api.DELETE("/books/:id/refs/:refid", bookid_param, app.ref_delete)

	g.Static("/app", "client")

	g.GET("/authtest", auth.AuthRequired(), func(ctx *gin.Context) {
		login := ctx.GetString(oidcauth.AuthUserKey)
		ctx.String(http.StatusOK, "Hi "+login)
	})

	g.Run("localhost:5555")
}
