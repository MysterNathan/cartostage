export interface User {
    id: number
    username: string
    first_name: string
    last_name: string
    email: string
    role: 'admin' | 'tutor' | 'student' | 'teacher' | 'enterprise'
    is_active: boolean
    created_at: string
    updated_at: string
    last_login?: string
}

export interface Tutor extends User {
    role: 'tutor'
}

export interface Teacher extends User {
    role: 'teacher'
}
export interface Student extends User {
    role: 'student'
}
// DTOs pour les requêtes
export interface CreateUserRequest {
    username: string
    first_name: string
    last_name: string
    email: string
    password: string
    role: 'admin' | 'tutor' | 'student' | 'teacher' | 'enterprise'
    entity_type?: string
    entity_id?: number
    phone?: string
    poste?: string
    departement?: string
}

export interface UpdateUserRequest {
    username?: string
    first_name?: string
    last_name?: string
    email?: string
    role?: string
    entity_type?: string
    entity_id?: number
    is_active?: boolean
}

export interface UpdateUserProfileRequest {
    phone?: string
    poste?: string
    departement?: string
    is_active?: boolean
}

export interface LoginRequest {
    username: string
    password: string
}

export interface LoginResponse {
    token: string
    user: User
    expires_at: number
}

export interface ChangePasswordRequest {
    old_password: string
    new_password: string
}

