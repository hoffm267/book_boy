import { useState, FormEvent, ChangeEvent } from 'react'
import type { AuthResponse } from '../types'

interface AuthProps {
  onLogin: (authData: AuthResponse) => void
  apiUrl: string
}

interface FormData {
  username: string
  email: string
  password: string
}

function Auth({ onLogin, apiUrl }: AuthProps) {
  const [isLogin, setIsLogin] = useState(true)
  const [formData, setFormData] = useState<FormData>({
    username: '',
    email: '',
    password: ''
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    const endpoint = isLogin ? '/auth/login' : '/auth/register'
    const body = isLogin
      ? { email: formData.email, password: formData.password }
      : formData

    try {
      const response = await fetch(`${apiUrl}${endpoint}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      })

      const data = await response.json().catch(() => null)

      if (!response.ok) {
        throw new Error(data?.error || 'Authentication failed')
      }

      onLogin(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Authentication failed')
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  const handleDemoLogin = async () => {
    setError('')
    setLoading(true)

    try {
      const response = await fetch(`${apiUrl}/auth/demo`, {
        method: 'POST',
      })

      const data = await response.json().catch(() => null)

      if (!response.ok) {
        throw new Error(data?.error || 'Demo account not available')
      }

      onLogin(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Demo account not available')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="card auth-card">
      <h2>{isLogin ? 'Login' : 'Register'}</h2>

      {error && <div className="error">{error}</div>}

      <form onSubmit={handleSubmit}>
        {!isLogin && (
          <div className="form-group">
            <label>Username</label>
            <input
              type="text"
              name="username"
              value={formData.username}
              onChange={handleChange}
              required={!isLogin}
            />
          </div>
        )}

        <div className="form-group">
          <label>Email</label>
          <input
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label>Password</label>
          <input
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            required
          />
        </div>

        <button type="submit" className="btn btn-primary" disabled={loading}>
          {loading ? 'Loading...' : (isLogin ? 'Login' : 'Register')}
        </button>
      </form>

      <div className="auth-divider">
        <span>or</span>
      </div>

      {isLogin && (
        <button
          type="button"
          onClick={handleDemoLogin}
          className="btn btn-demo"
          disabled={loading}
        >
          {loading ? 'Loading...' : 'Try Demo Account — No sign-up needed'}
        </button>
      )}

      <p className="auth-switch">
        {isLogin ? "Don't have an account? " : 'Already have an account? '}
        <button onClick={() => setIsLogin(!isLogin)} className="link-btn">
          {isLogin ? 'Register' : 'Login'}
        </button>
      </p>
    </div>
  )
}

export default Auth
