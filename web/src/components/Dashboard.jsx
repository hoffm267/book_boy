import Progress from './Progress'

function Dashboard({ user, token, onLogout, apiUrl }) {
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

            <Progress token={token} apiUrl={apiUrl} userId={user.id} />
        </>
    )
}

export default Dashboard
