export class UnauthorizedError extends Error {
  constructor() {
    super('Unauthorized - Session expired')
    this.name = 'UnauthorizedError'
  }
}

export async function fetchWithAuth(url, options = {}, onUnauthorized) {
  const response = await fetch(url, options)

  if (response.status === 401) {
    console.warn('Received 401 Unauthorized - logging out')
    if (onUnauthorized) {
      onUnauthorized()
    }
    throw new UnauthorizedError()
  }

  return response
}
