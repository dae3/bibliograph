import React from 'react';
import { useState, useRef } from 'react';
import { useQueryClient, useQuery, useMutation } from '@tanstack/react-query'
import Editor  from './Editor.js'
import Graph  from './Graph.js'

export default function App({apibase}) {
  const queryClient = useQueryClient()
  const { isPending, isError, data, error } = useQuery({
    queryKey: ['books'],
    queryFn: async () => {
      const r = await fetch(`${apibase}/books`)
      if (r.ok) {
        return r.json()
      } else {
        throw new Error(r.statusText)
      }
    }
  })
  const deleteBook = useMutation({
    mutationFn: (id) => fetch(`${apibase}/books/${id}`, { method: 'DELETE' }),
    onSuccess: () => { queryClient.invalidateQueries(['books']) }
  })

  const addBook = useMutation({
    mutationFn: ({author, title}) => fetch(`${apibase}/books`,
      {
        method: 'POST',
        headers: new Headers({'Content-Type':'application/json'}),
        body: JSON.stringify({ 'author': author, 'title': title })
      }
    ),
    onSuccess: () => { queryClient.invalidateQueries(['books']) }
  })

  const addRef = useMutation({
    mutationFn: ({ sourceid, refid }) => fetch(`${apibase}/books/${sourceid}/refs`,
      {
        method: 'POST',
        headers: new Headers({'Content-Type':'application/json'}),
        body: JSON.stringify({ refs: data.find(b => b.id == sourceid).references.concat(refid) })
      }
    ),
    onSuccess: () => { queryClient.invalidateQueries(['books']) }
  })

  const delRef = useMutation({
    mutationFn: ({ sourceid, refid }) => fetch(`${apibase}/books/${sourceid}/refs/${refid}`, { method: 'DELETE'}),
    onSuccess: () => { queryClient.invalidateQueries(['books']) }
  })

  if (isPending) {
    return(<p>loading</p>)
  }

  if (isError) {
    return(<p>{error}</p>)
  }

  return(
    <>
      <Graph books={data} />
      <Editor books={data} deleteBook={deleteBook.mutate} addBook={addBook.mutate} addRef={addRef.mutate} delRef={delRef.mutate}/>
    </>
  )
}
