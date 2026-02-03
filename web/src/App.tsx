import { useState, useEffect } from 'react'
import Auth from './components/Auth'
import Dashboard from './components/Dashboard'
import type { User, AuthResponse } from './types'

const API_URL = import.meta.env.VITE_API_URL || '/api'

function App() {
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'))
  const [user, setUser] = useState<User | null>(null)

  useEffect(() => {
    const storedUser = localStorage.getItem('user')
    if (storedUser) {
      setUser(JSON.parse(storedUser))
    }
  }, [])

  const handleLogin = (authData: AuthResponse) => {
    setToken(authData.token)
    setUser(authData.user)
    localStorage.setItem('token', authData.token)
    localStorage.setItem('user', JSON.stringify(authData.user))
  }

  const handleLogout = () => {
    setToken(null)
    setUser(null)
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return (
    <div className="container">
      {!token ? (
        <Auth onLogin={handleLogin} apiUrl={API_URL} />
      ) : (
        <Dashboard
          user={user}
          token={token}
          onLogout={handleLogout}
          apiUrl={API_URL}
        />
      )}
    </div>
  )
}

export default App
