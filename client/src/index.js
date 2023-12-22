import './main.css';
import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App.js'
import {
  QueryClient,
  QueryClientProvider,
  useQuery,
} from '@tanstack/react-query'

const queryClient = new QueryClient()

// graph(data, document.getElementById('graph'))
createRoot(document.getElementById('app')).render(
  <>
    <QueryClientProvider client={queryClient}>
      <App apibase="http://localhost:5555/api/v1" />
    </QueryClientProvider>
  </>
)
