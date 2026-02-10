class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

export function createApi(baseUrl: string, token: string, onUnauthorized: () => void) {
  async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const headers: Record<string, string> = {
      'Authorization': `Bearer ${token}`,
      ...options.headers as Record<string, string>,
    }

    if (options.body) {
      headers['Content-Type'] = 'application/json'
    }

    const response = await fetch(`${baseUrl}${path}`, { ...options, headers })

    if (response.status === 401) {
      onUnauthorized()
      throw new ApiError('Unauthorized', 401)
    }

    if (!response.ok) {
      const data = await response.json().catch(() => ({}))
      throw new ApiError(data.error || 'Request failed', response.status)
    }

    const text = await response.text()
    return text ? JSON.parse(text) : (undefined as T)
  }

  return {
    get: <T>(path: string) => request<T>(path),
    post: <T>(path: string, body?: unknown) => request<T>(path, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    }),
    patch: <T>(path: string, body: unknown) => request<T>(path, {
      method: 'PATCH',
      body: JSON.stringify(body),
    }),
    delete: (path: string) => request<void>(path, { method: 'DELETE' }),
  }
}

export type Api = ReturnType<typeof createApi>
