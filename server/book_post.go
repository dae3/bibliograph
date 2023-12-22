package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *App) book_post(ctx *gin.Context) {
	b := new(APIBook)
	if err := ctx.BindJSON(b); err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	newbook, err := app.db.Book.Create().SetAuthor(b.Author).SetTitle(b.Title).Save(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else {
		ctx.JSON(http.StatusOK, ApiBookFromBook(newbook))
	}
}
