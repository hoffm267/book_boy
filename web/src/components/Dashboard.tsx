import Progress from './Progress'
import type { User } from '../types'

interface DashboardProps {
  user: User | null
  token: string
  onLogout: () => void
  apiUrl: string
}

function Dashboard({ user, token, onLogout, apiUrl }: DashboardProps) {
  if (!user) return null

  return (
    <>
      <div className="header">
        <div>
          <h1>Book Boy</h1>
          <p style={{ color: '#7f8c8d', marginTop: '5px' }}>Welcome, {user.username}</p>
        </div>
        <button onClick={onLogout} className="btn btn-secondary">
          Logout
        </button>
      </div>

      <Progress token={token} apiUrl={apiUrl} userId={user.id} onLogout={onLogout} />
    </>
  )
}

export default Dashboard
