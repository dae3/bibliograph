package main

import (
	"bibliograph/api/ent"
	"bibliograph/api/ent/book"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) ref_post(ctx *gin.Context) {
	var body struct {
		Refs []int `json:"refs"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	source, err := app.db.Book.Query().WithReferences().Where(book.IDEQ(ctx.GetInt(ParamBookId))).Only(ctx.Request.Context())
	if ent.IsNotFound(err) {
		ctx.Status(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	} else {
		for _, refid := range body.Refs {
			_, err := app.db.Book.Get(ctx.Request.Context(), refid)
			if ent.IsNotFound(err) {
				ctx.Status(http.StatusUnprocessableEntity)
				return
			} else if err != nil {
				ctx.Error(err)
				ctx.Status(http.StatusInternalServerError)
				return
			}
		}
		source, err = source.Update().AddReferenceIDs(body.Refs...).Save(ctx.Request.Context())
		if err != nil {
			ctx.Error(err)
			ctx.Status(http.StatusInternalServerError)
			return
		} else {
			ctx.JSON(http.StatusOK, ApiBookFromBook(source))
		}
	}
}
