package main

import (
	"bibliograph/api/ent"
	"bibliograph/api/ent/book"
	"net/http"
)

func (app *App) book_delete(w http.ResponseWriter, r *http.Request) {
	bookid, err := GetIntParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	c, err := app.db.Book.Delete().Where(book.IDEQ(bookid)).Exec(r.Context())
	if c == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (app *App) ref_delete(w http.ResponseWriter, r *http.Request) {
	bookid, err := GetIntParam(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	refid, err := GetIntParam(r, "refid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := app.db.Book.Get(r.Context(), bookid)
	if ent.IsNotFound(err) {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = book.Update().RemoveReferenceIDs(int(refid)).Exec(r.Context())
	if ent.IsNotFound(err) {
		http.NotFound(w, r)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
