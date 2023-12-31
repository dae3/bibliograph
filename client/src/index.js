import './main.css';
import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.js'
import { QueryClient, QueryClientProvider, useQuery } from '@tanstack/react-query'

const queryClient = new QueryClient()
const backend = { // will be substituted by webpack_U
  base: BASE_URL,
  api: API_BASE,
  auth: AUTH_BASE
}

console.log(backend)

createRoot(document.getElementById('app')).render(
  <>
    <QueryClientProvider client={queryClient}>
      <App backend={backend}/>
    </QueryClientProvider>
  </>
)
