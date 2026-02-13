import Progress from './Progress'
import { useAuth } from '../AuthContext'

function Dashboard() {
  const { user, onLogout } = useAuth()

  return (
    <>
      <div className="header">
        <div>
          <h1>Book Boy</h1>
          <p className="subtitle">Welcome, {user.username}</p>
        </div>
        <button onClick={onLogout} className="btn btn-secondary">
          Logout
        </button>
      </div>

      <Progress />
    </>
  )
}

export default Dashboard
