import React from 'react';
import { useState, useRef } from 'react';

export default function Editor({ loggedIn, books, bookid, deleteBook, addBook, updateBook, addRef, delRef, removeSelection }) {
  const withrefs = books.map(b => enrichReferences(books, b))
  const book = withrefs.find(b => b.id == bookid)
  const availrefs = books // not self or existing ref
    .filter(b => b.id != book?.id)
    .filter(b => !book?.references.map(r => r.id).includes(b.id))

  const [ author, setAuthor ] = useState(book ? book.author : '')
  const [ title, setTitle ] = useState(book ? book.title : '')
  const newref = useRef(null)

  if (!loggedIn) { return("") }
  else { return(
    <div className="inline-block space-x-1 max-wd-md">
      <label className="w-1/3">Title<input className="ml-0.5 border" type="text" value={title} onChange={e => setTitle(e.target.value)} /></label>
      <label className="w-1/3">Author<input className="ml-0.5 border" type="text" value={author} onChange={e => setAuthor(e.target.value)} /></label>
      {
        book
          ? <>
            <input type="button" className="bg-slate-200 rounded-sm p-0.5" value="Update" onClick={() => {updateBook({id: book.id, author: author, title: title})} } />
            <input type="button" className="bg-slate-200 rounded-sm p-0.5" value="Delete" onClick={() => {deleteBook(book.id)}} />
            <input className="bg-slate-200 rounded-sm" type="button" value="Cancel" onClick={removeSelection} />
            <select name="addref" id="addref" ref={newref}>{availrefs?.map(b => <option key={b.id} value={b.id}><BookDisplay book={b} /></option>)}</select>
            <input className="bg-slate-200 rounded-sm p-0.5" type="button" value="Add reference" onClick={() => { addRef({ sourceid: book.id, refid: parseInt(newref.current.value)}) }} />
            <References book={book} />
          </>
          : <input className="bg-slate-200 rounded-sm p-0.5" type="button" value="Add" onClick={() => {addBook({author: author, title: title}); setAuthor(''); setTitle('')}}/>
      }
    </div>
  )}
}

function References({ book }) {
  return(
    <ul>
      { book.references.map(r =>
      <li key={r.id} className="max-wi">
        <BookDisplay book={r}/>
        <input className="bg-slate-200 rounded-sm p-0.5" type="button" value="Delete reference" onClick={() => delRef({sourceid: book.id, refid: r.id})} />
      </li>
      )}
    </ul>
  )
}

function BookDisplay({ book }) {
  return(`${book.title}: ${book.author}`)
}

// utility functions
function enrichReferences(books, book) {
  const ebook = {...book}
  ebook.references = book.references.map(id => books.filter(b => b.id == id)[0])
  return(ebook)
}

// vim: ft=javascript.jsx
