import { createContext, useContext, useMemo } from 'react'
import { createApi } from './utils/api'
import type { Api } from './utils/api'
import type { User } from './types'

interface AuthContextType {
  token: string
  user: User
  apiUrl: string
  api: Api
  onLogout: () => void
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children, token, user, apiUrl, onLogout }: {
  children: React.ReactNode
  token: string
  user: User
  apiUrl: string
  onLogout: () => void
}) {
  const api = useMemo(() => createApi(apiUrl, token, onLogout), [apiUrl, token, onLogout])

  return (
    <AuthContext.Provider value={{ token, user, apiUrl, api, onLogout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
