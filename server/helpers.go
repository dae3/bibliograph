package main

import (
	"bibliograph/api/ent"
	_ "bibliograph/api/ent/book"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type APIBook struct {
	Id         int    `json:"id"`
	Author     string `json:"author" validate:"required"`
	Title      string `json:"title" validate:"required"`
	References []int  `json:"references"`
}

func ApiBookFromBook(book *ent.Book) (apibook APIBook) {
	apibook = APIBook{book.ID, book.Author, book.Title, make([]int, len(book.Edges.References))}
	for k, v := range book.Edges.References {
		apibook.References[k] = v.ID
	}
	return
}

func GetIntParam(r *http.Request, param string) (bookid int, err error) {
	vars := mux.Vars(r)
	bookidparam, ok := vars[param]
	if !ok {
		err = errors.Join(err, errors.New(fmt.Sprintf("Parameter %s not found", param)))
		return
	}
	bookid64, err := (strconv.ParseInt(bookidparam, 0, 0))
	if err != nil {
		err = errors.Join(err, errors.New("Can't parse book id as integer"))
	} else {
		bookid = int(bookid64)
	}
	return
}

// func ValidateContentTypeMiddleware(ctx *gin.Context) error {
// 	if ctx.Method() == fiber.MethodPost {
// 		if ctx.Is("json") {
// 			return ctx.Next()
// 		} else {
// 			return fiber.ErrUnsupportedMediaType
// 		}
// 	} else {
// 		return ctx.Next()
// 	}
// }
