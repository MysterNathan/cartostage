'use client';

import { useState, useEffect, createContext, useContext } from 'react';
import { authApi, LoginRequest } from '@/lib/api/authApi';

interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => void;
  checkAuth: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const useAuthState = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  const checkAuth = () => {
    const authenticated = authApi.isAuthenticated();
    setIsAuthenticated(authenticated);
    setIsLoading(false);
  };

  const login = async (credentials: LoginRequest) => {
    try {
      await authApi.login(credentials);
      setIsAuthenticated(true);
    } catch (error) {
      setIsAuthenticated(false);
      throw error;
    }
  };

  const logout = () => {
    authApi.logout();
    setIsAuthenticated(false);
  };

  useEffect(() => {
    checkAuth();
  }, []);

  return {
    isAuthenticated,
    isLoading,
    login,
    logout,
    checkAuth,
  };
};
