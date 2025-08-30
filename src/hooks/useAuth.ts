'use client'
import { useState, useEffect } from 'react'
import { getStoredAuth, storeAuth, clearAuth, validateCredentials } from '@/lib/auth'
import type { User } from '@/lib/auth'

export function useAuth() {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const storedAuth = getStoredAuth()
    setUser(storedAuth)
    setLoading(false)
  }, [])

  const login = (username: string, password: string): boolean => {
    if (validateCredentials(username, password)) {
      const user: User = { username, isAuthenticated: true }
      setUser(user)
      storeAuth(user)
      return true
    }
    return false
  }

  const logout = () => {
    setUser(null)
    clearAuth()
  }

  return {
    user,
    loading,
    isAuthenticated: user?.isAuthenticated || false,
    login,
    logout
  }
}
