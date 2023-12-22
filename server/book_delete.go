package main

import (
	"bibliograph/api/ent"
	"bibliograph/api/ent/book"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *App) book_delete(ctx *gin.Context) {
	c, err := app.db.Book.Delete().Where(book.IDEQ(ctx.GetInt(ParamBookId))).Exec(ctx.Request.Context())
	if c == 0 {
		ctx.Status(http.StatusNotFound)
	} else if err != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Error(err)
	} else {
		ctx.Status(http.StatusOK)
	}
}

func (app *App) ref_delete(ctx *gin.Context) {
	refid, err := strconv.ParseInt(ctx.Param("refid"), 0, 0)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	book, err := app.db.Book.Get(ctx.Request.Context(), ctx.GetInt(ParamBookId))
	if ent.IsNotFound(err) {
		ctx.Status(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Error(err)
		return
	}

	err = book.Update().RemoveReferenceIDs(int(refid)).Exec(ctx.Request.Context())
	if ent.IsNotFound(err) {
		ctx.Status(http.StatusNotFound)
	} else if err != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Error(err)
	} else {
		ctx.Status(http.StatusOK)
	}
}
