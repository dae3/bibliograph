package main

import (
	"bibliograph/api/ent"
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *App) book_get(w http.ResponseWriter, r *http.Request) {

	bookid, err := GetIntParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	book, err := app.db.Book.Get(r.Context(), bookid)
	if ent.IsNotFound(err) {
		w.WriteHeader(http.StatusNotFound)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s\n", err.Error())
	} else {
		fmt.Printf("%+v", book)
		w.WriteHeader(http.StatusOK)
		j := json.NewEncoder(w)
		j.Encode(ApiBookFromBook(book))
	}
}

func (app *App) books_get(w http.ResponseWriter, r *http.Request) {
	if books, err := app.db.Book.Query().WithReferences().All(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		apibooks := make([]APIBook, len(books))
		for k, v := range books {
			apibooks[k] = ApiBookFromBook(v)
		}
		w.WriteHeader(http.StatusOK)
		j := json.NewEncoder(w)
		j.Encode(apibooks)
	}
}
