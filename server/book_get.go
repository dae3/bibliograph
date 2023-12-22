package main

import (
	"bibliograph/api/ent"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) book_get(ctx *gin.Context) {
	book, err := app.db.Book.Get(ctx.Request.Context(), ctx.GetInt(ParamBookId))
	if ent.IsNotFound(err) {
		ctx.Status(http.StatusNotFound)
	} else if err != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, ApiBookFromBook(book))
	}
}

func (app *App) books_get(ctx *gin.Context) {
	if books, err := app.db.Book.Query().WithReferences().All(ctx.Request.Context()); err != nil {
		ctx.Status(http.StatusInternalServerError)
	} else {
		apibooks := make([]APIBook, len(books))
		for k, v := range books {
			apibooks[k] = ApiBookFromBook(v)
		}
		ctx.JSON(http.StatusOK, apibooks)
	}
}
