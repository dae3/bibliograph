package main

import (
	"bibliograph/api/ent"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (app *App) book_post(w http.ResponseWriter, r *http.Request) {
	apibook, err := parsePost(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := GetIntParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := app.db.Book.Get(r.Context(), id)
	if ent.IsNotFound(err) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	book, err = book.Update().SetAuthor(apibook.Author).SetTitle(apibook.Title).Save(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(ApiBookFromBook(book))
	}
	return
}

func (app *App) book_post_new(w http.ResponseWriter, r *http.Request) {
	apibook, err := parsePost(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := app.db.Book.Create().SetAuthor(apibook.Author).SetTitle(apibook.Title).Save(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(ApiBookFromBook(book))
	}
}

func parsePost(r *http.Request) (apibook *APIBook, err error) {
	body := make([]byte, r.ContentLength, r.ContentLength)
	_, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		return
	}

	// accept liberally but validate the real fields
	apibook = new(APIBook)
	err = json.Unmarshal(body, apibook)
	if err != nil {
		return
	}
	if apibook.Author == "" || apibook.Title == "" {
		err = errors.New("Missing required fields")
	}

	return
}
