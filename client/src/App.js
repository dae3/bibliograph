import React from 'react';
import { useState, useRef } from 'react';
import { useQueryClient, useQuery, useMutation } from '@tanstack/react-query'
import Editor from './Editor.js'
import Graph from './Graph.js'
import Auth from './Auth.js'

export default function App({ backend }) {
  const [ selectedNode, setSelectedNode ] = useState(null)
  const [ authState, setAuthState ] = useState({ username: null, isLoggedIn: false })
  const queryClient = useQueryClient()
  const csrfToken = useRef(null)

  const req = function(url, method = 'GET', body = null, headers = {}) {
    if (csrfToken.current == null) {
      fetch(`${backend.base}${backend.api}token`, { credentials: 'include' }).then(r => r.text()).then(t => csrfToken.current = t)
    }
    if (method != 'GET') {
      headers = Object.assign(headers, {'X-CSRF-Token': csrfToken.current})
    }
    return new Request(url, { method: method, body: body, credentials: "include", headers: new Headers(headers) })
  }

  const { isPending, isError, data, error } = useQuery({
    queryKey: ['books'],
    queryFn: async () => {
      const r = await fetch(req(`${backend.base}${backend.api}books`))
      if (r.ok) {
        return r.json()
      } else {
        throw new Error(r.statusText)
      }
    }
  })

  const delBook = useMutation({
    mutationFn: (id) => fetch(req(`${backend.base}${backend.api}books/${id}`, 'DELETE')),
    onSuccess: () => { queryClient.invalidateQueries(['books']); setSelectedNode(null) }
  })

  const addBook = useMutation({
    mutationFn: ({author, title}) => fetch(req(
      `${backend.base}${backend.api}books`,
      'POST',
      JSON.stringify({ 'author': author, 'title': title }),
      { 'Content-Type': 'application/json' }
    )),
    onSuccess: () => { queryClient.invalidateQueries(['books']); setSelectedNode(null) }
  })

  const updateBook = useMutation({
    mutationFn: ({id, author, title}) => fetch(req(
      `${backend.base}${backend.api}books/${id}`,
      'POST',
      JSON.stringify({ 'author': author, 'title': title }),
      {'Content-Type':'application/json'},
    )),
    onSuccess: () => { queryClient.invalidateQueries(['books']); setSelectedNode(null) }
  })

  const addRef = useMutation({
    mutationFn: ({ sourceid, refid }) => fetch(req(
      `${backend.base}${backend.api}books/${sourceid}/refs`,
      'POST',
      JSON.stringify({ refs: data.find(b => b.id == sourceid).references.concat(refid) }),
      {'Content-Type':'application/json'},
    )),
    onSuccess: () => { queryClient.invalidateQueries(['books']) }
  })

  const delRef = useMutation({
    mutationFn: ({ sourceid, refid }) => fetch(req(`${backend.base}${backend.api}books/${sourceid}/refs/${refid}`, 'DELETE')),
    onSuccess: () => { queryClient.invalidateQueries(['books']) }
  })

  if (isPending) {
    return(<p>Loading...</p>)
  }

  if (isError) {
    return(<p>{error}</p>)
  }

  return(
    <>
      <header className="bg-blue-500">
        <h1 className="text-lg">Bibliograph</h1>
        <nav>
          <Auth backend={`${backend.base}${backend.auth}`} apireq={req} setAuthState={setAuthState} authState={authState}/>
        </nav>
      </header>
      <Editor loggedIn={authState.loggedIn} books={data} key={selectedNode ? selectedNode : ''} bookid={selectedNode} addBook={addBook.mutate} updateBook={updateBook.mutate} deleteBook={delBook.mutate} addRef={addRef.mutate} delRef={delRef.mutate} removeSelection={() => setSelectedNode(null)}/>
      <Graph books={data} nodeSelected={setSelectedNode} />
    </>
  )
}
