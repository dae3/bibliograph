import React from 'react';

export default function Editor({ data }) {
  return(<ul>{data.map(book => <Book book={book} />)}</ul>)
}

function Book({ book }) {
  var reflist

  if (Object.hasOwn(book, 'references')) {
    const refs = book.references.map(ref => <li key={ref}>aref</li>)
    reflist = <ul>{refs}</ul>
  } else {
    reflist = null
  }

  return(
    <>
      <li key={book.id}>{book.title} - {book.author}</li>
      {reflist}
    </>
  )
}
