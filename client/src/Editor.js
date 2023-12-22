import React from 'react';
import { useState, useRef } from 'react';

export default function Editor({ books, deleteBook, addBook, addRef, delRef }) {
  const withrefs = books.map(book => enrichReferences(books, book))
  const newauthor = useRef(null)
  const newtitle = useRef(null)
  const [ addingRef, setAddingRef ] = useState({ active: false, bookid: null })
  const newref = useRef(null)

  function addRefToggle(bookid) {
    if (addingRef.active) {
      setAddingRef({ active: false, bookid: null})
    } else {
      setAddingRef({ active: true, bookid: bookid})
    }
  }

  function addRefControl(bookid) {
    return(
      <>
        <label htmlFor="addref">New reference</label>
        <select name="addref" id="addref" ref={newref}>
          {books.filter(book => book.id != bookid).map(book => <option value={book.id}><BookDisplay book={book} /></option>)}
        </select>
        <input type="submit" value="add" onClick={() => {
          addRef({ sourceid: bookid, refid: parseInt(newref.current.value)})
          addRefToggle(bookid)
        }}
        />
      </>
    )
  }

  return(
    <>
      <ul id="editor">
        {withrefs.map(book =>
            <>
              <Book book={book} delRef={delRef} />
              <span className="book" onClick={() => {deleteBook(book.id)}}> X </span>
              {addingRef.active && addingRef.bookid == book.id ? <span>{addRefControl(book.id)}</span> : <span onClick={() => addRefToggle(book.id)}>R</span>}
            </>
        )}
      </ul>
      <label htmlFor="newtitle">New title</label>
      <input type="text" id="newtitle" name="newtitle" ref={newtitle} />
      <label htmlFor="newauthor">New author</label>
      <input type="text" id="newauthor" name="newauthor" ref={newauthor} />
      <input type="submit" value="Add" onClick={() => {addBook({author: newauthor.current.value, title: newtitle.current.value})}} />
    </>
  )
}

// supporting components
function Book({ book, delRef }) {
  const [ isExpanded, setIsExpanded ] = useState(false)

  return(
    <li
      className={[ book.references.length > 0 ? "hasrefs" : "", isExpanded ? "expanded" : "" ].join(' ')}
      key={book.id}
      onClick={() => setIsExpanded(!isExpanded)}
    >
      <BookDisplay book={book} />
      {isExpanded ? <ul>{ book.references.map(r => <li key={book.id}><BookDisplay book={r} /><span onClick={() => delRef({sourceid: book.id, refid: r.id})}> X </span></li>) }</ul> : ""}
    </li>
  )
}

function BookDisplay({ book }) {
  return(<><span className="title">{ book.title }</span> <span className="author">{ book.author }</span></>)
}

// utility functions
function enrichReferences(books, book) {
  const ebook = {...book}
  ebook.references = book.references.map(id => books.filter(b => b.id == id)[0])
  return(ebook)
}

// vim: ft=javascript.jsx
