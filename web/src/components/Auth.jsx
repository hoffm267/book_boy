import { useState } from 'react'

function Auth({ onLogin, apiUrl }) {
    const [isLogin, setIsLogin] = useState(true)
    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: ''
    })
    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)

    const handleSubmit = async (e) => {
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

            const data = await response.json()

            if (!response.ok) {
                throw new Error(data.error || 'Authentication failed')
            }

            onLogin(data)
        } catch (err) {
            setError(err.message)
        } finally {
            setLoading(false)
        }
    }

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value })
    }

    const handleDemoLogin = async () => {
        setError('')
        setLoading(true)

        try {
            const response = await fetch(`${apiUrl}/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    email: 'demo@bookboy.app',
                    password: 'Demo123!'
                })
            })

            const data = await response.json()

            if (!response.ok) {
                throw new Error(data.error || 'Demo login failed')
            }

            onLogin(data)
        } catch (err) {
            setError(err.message)
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="card" style={{ maxWidth: '400px', margin: '100px auto' }}>
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

        {isLogin && (
            <button
            type="button"
            onClick={handleDemoLogin}
            className="btn"
            disabled={loading}
            style={{
                width: '100%',
                    marginTop: '10px',
                    backgroundColor: '#95a5a6',
                    color: 'white',
                    border: 'none'
            }}
            >
            {loading ? 'Loading...' : 'Try Demo Account'}
            </button>
        )}

        <p style={{ marginTop: '15px', textAlign: 'center' }}>
        {isLogin ? "Don't have an account? " : 'Already have an account? '}
        <button
        onClick={() => setIsLogin(!isLogin)}
        style={{ background: 'none', border: 'none', color: '#3498db', cursor: 'pointer' }}
        >
        {isLogin ? 'Register' : 'Login'}
        </button>
        </p>
        </div>
    )
}

export default Auth
