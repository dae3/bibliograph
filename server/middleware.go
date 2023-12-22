package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	ParamBookId = "bookid"
)

// bookid_param attempts to parse an integer path parameter named "id".
// If successful the int value of the parameter is added to the context
// with the key "bookid", otherwise http.StatusBadRequest is returned
func bookid_param(ctx *gin.Context) {
	idparam := ctx.Param("id")
	if idparam == "" {
		ctx.Status(http.StatusBadRequest)
		ctx.Abort()
	} else {
		id, err := strconv.ParseInt(idparam, 0, 0)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			ctx.Error(err)
			ctx.Abort()
		} else {
			ctx.Set(ParamBookId, int(id))
			ctx.Next()
		}
	}
}
