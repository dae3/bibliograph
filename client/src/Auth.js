import React from 'react';
import { useState, useEffect } from 'react';

export default function Auth({ backend, apireq, authState, setAuthState }) {
  useEffect(() => {
    fetch(apireq(`${backend}status`))
      .then(r => {
        if (r.status == 200 && !authState.loggedIn) {
          r.json().then(j => setAuthState({username: j.email, loggedIn: true}))
        }
      })
  })

  return(
    <ul className="space-x-2">
      { !authState.loggedIn && <li className="inline-block"><a href={`${backend}/login`}>Login</a></li> }
      { authState.loggedIn && <li className="inline-block"><a href={`${backend}/logout`}>Logout</a></li> }
      <span>{authState?.username}</span>
    </ul>
  )
}
