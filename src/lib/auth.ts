export interface User {
  username: string
  isAuthenticated: boolean
}

export const AUTH_CONFIG = {
  username: 'admin',
  password: 'lycee2024'
}

export function validateCredentials(username: string, password: string): boolean {
  return username === AUTH_CONFIG.username && password === AUTH_CONFIG.password
}

export function getStoredAuth(): User | null {
  if (typeof window === 'undefined') return null
  
  // Utiliser localStorage au lieu de sessionStorage pour persister
  const stored = localStorage.getItem('auth_session') // Utiliser la même clé
  if (!stored) return null
  
  try {
    return JSON.parse(stored)
  } catch {
    return null
  }
}

export function storeAuth(user: User): void {
  if (typeof window === 'undefined') return
  localStorage.setItem('auth_session', JSON.stringify(user)) // Utiliser la même clé
}

export function clearAuth(): void {
  if (typeof window === 'undefined') return
  localStorage.removeItem('auth_session') // Utiliser la même clé
}

export function isAuthenticated(): boolean {
  const user = getStoredAuth()
  return user?.isAuthenticated === true
}
