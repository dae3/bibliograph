package main

import (
	"bibliograph/api/ent"
	"bibliograph/api/ent/book"
	"encoding/json"
	"io"
	"net/http"
)

func (app *App) ref_post(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	var refs struct {
		Refs []int `json:"refs"`
	}
	nbytes, err := r.Body.Read(body)
	if nbytes < int(r.ContentLength) || (err != nil && err != io.EOF) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &refs)
	if err != nil || len(refs.Refs) == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookid, err := GetIntParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	source, err := app.db.Book.Query().WithReferences().Where(book.IDEQ(bookid)).Only(r.Context())
	if ent.IsNotFound(err) {
		http.Error(w, "Book id not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		for _, refid := range refs.Refs {
			_, err := app.db.Book.Get(r.Context(), refid)
			if ent.IsNotFound(err) {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		source, err = source.Update().AddReferenceIDs(refs.Refs...).Save(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			json.NewEncoder(w).Encode(ApiBookFromBook(source))
		}
	}
}
