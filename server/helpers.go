package main

import (
	"bibliograph/api/ent"
	_ "bibliograph/api/ent/book"
)

type APIBook struct {
	Id         int    `json:"id"`
	Author     string `json:"author"`
	Title      string `json:"title"`
	References []int  `json:"references"`
}

func ApiBookFromBook(book *ent.Book) (apibook APIBook) {
	apibook = APIBook{book.ID, book.Author, book.Title, make([]int, len(book.Edges.References))}
	for k, v := range book.Edges.References {
		apibook.References[k] = v.ID
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
