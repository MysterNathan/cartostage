const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  expires_at: string;
}

export interface CreateUserRequest {
  username: string;
  password: string;
  role?: string;
}

export interface User {
  id: number;
  username: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface UpdateUserRequest {
  username: string;
  role: string;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

class AuthApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'AuthApiError';
  }
}

// Utilitaire pour récupérer le token depuis le localStorage
const getAuthToken = (): string | null => {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem('authToken');
};

// Utilitaire pour les headers avec authentification
const getAuthHeaders = (): HeadersInit => {
  const token = getAuthToken();
  return {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': `Bearer ${token}` }),
  };
};

// Utilitaire pour gérer les réponses API
const handleApiResponse = async <T>(response: Response): Promise<T> => {
  if (!response.ok) {
    const errorText = await response.text();
    throw new AuthApiError(response.status, errorText || `HTTP ${response.status}`);
  }

  // Si la réponse est vide (204 No Content), retourner null
  if (response.status === 204) {
    return null as T;
  }

  return response.json();
};

export const authApi = {
  // Connexion
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response = await fetch(`${API_BASE_URL}/api/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials),
    });

    const data = await handleApiResponse<LoginResponse>(response);

    // Stocker le token en localStorage après une connexion réussie
    if (data.token) {
      localStorage.setItem('authToken', data.token);
      localStorage.setItem('authTokenExpiry', data.expires_at);
    }

    return data;
  },

  // Déconnexion (côté client uniquement)
  logout: (): void => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('authToken');
      localStorage.removeItem('authTokenExpiry');
    }
  },

  // Vérifier si l'utilisateur est connecté
  isAuthenticated: (): boolean => {
    if (typeof window === 'undefined') return false;

    const token = localStorage.getItem('authToken');
    const expiry = localStorage.getItem('authTokenExpiry');

    if (!token || !expiry) return false;

    // Vérifier si le token a expiré
    const expiryDate = new Date(expiry);
    const now = new Date();

    if (now >= expiryDate) {
      // Token expiré, le supprimer
      localStorage.removeItem('authToken');
      localStorage.removeItem('authTokenExpiry');
      return false;
    }

    return true;
  },

  // Récupérer le token actuel
  getToken: getAuthToken,

  // === Routes Admin (nécessitent une authentification) ===

  // Récupérer tous les utilisateurs
  getUsers: async (): Promise<User[]> => {
    const response = await fetch(`${API_BASE_URL}/admin/auth/users`, {
      method: 'GET',
      headers: getAuthHeaders(),
    });

    return handleApiResponse<User[]>(response);
  },

  // Créer un nouvel utilisateur
  createUser: async (userData: CreateUserRequest): Promise<User> => {
    const response = await fetch(`${API_BASE_URL}/admin/auth/users`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(userData),
    });

    return handleApiResponse<User>(response);
  },

  // Mettre à jour un utilisateur
  updateUser: async (userId: number, userData: UpdateUserRequest): Promise<User> => {
    const response = await fetch(`${API_BASE_URL}/admin/auth/users/${userId}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(userData),
    });

    return handleApiResponse<User>(response);
  },

  // Supprimer un utilisateur
  deleteUser: async (userId: number): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/admin/auth/users/${userId}`, {
      method: 'DELETE',
      headers: getAuthHeaders(),
    });

    return handleApiResponse<void>(response);
  },

  // Changer le mot de passe
  changePassword: async (userId: number, passwordData: ChangePasswordRequest): Promise<void> => {
    const response = await fetch(`${API_BASE_URL}/admin/auth/users/${userId}/password`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(passwordData),
    });

    return handleApiResponse<void>(response);
  },
};

export default authApi;
