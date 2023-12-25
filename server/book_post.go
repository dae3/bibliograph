package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (app *App) book_post(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength, r.ContentLength)
	nbytes, err := r.Body.Read(body)
	if nbytes < int(r.ContentLength) || (err != nil && err != io.EOF) {
		http.Error(w, fmt.Sprintf("Error reading body: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// accept liberally but validate the real fields
	newapibook := new(APIBook)
	if err := json.Unmarshal(body, newapibook); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if newapibook.Author == "" || newapibook.Title == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	newbook, err := app.db.Book.Create().SetAuthor(newapibook.Author).SetTitle(newapibook.Title).Save(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(ApiBookFromBook(newbook))
	}
}
