// lib/enterpriseApi.ts
import type { EnterpriseData } from '@/types/enterprise'
import type { Stage } from '@/types/stage'
import type { Tutor } from '@/types/tutor'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

// Helper pour gérer les erreurs d'API
async function handleApiResponse<T>(response: Response): Promise<T> {
    if (!response.ok) {
        const errorText = await response.text()

        if (response.status === 401) {
            localStorage.removeItem('authToken')
            throw new Error('Session expirée - Veuillez vous reconnecter')
        }

        if (response.status === 403) {
            throw new Error('Accès non autorisé pour cette action')
        }

        if (response.status === 404) {
            throw new Error('Ressource non trouvée')
        }

        let errorMessage = 'Erreur serveur'
        try {
            const errorData = JSON.parse(errorText)
            errorMessage = errorData.message || errorMessage
        } catch {
            errorMessage = errorText || errorMessage
        }

        throw new Error(errorMessage)
    }

    return response.json()
}

// Helper pour les headers d'authentification
function getAuthHeaders(): HeadersInit {
    const token = localStorage.getItem('authToken')
    if (!token) {
        throw new Error('Token d\'authentification manquant')
    }

    return {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
    }
}

export async function getMyEnterpriseData(): Promise<EnterpriseData> {
    const response = await fetch(`${API_URL}/api/enterprises/me`, {
        headers: getAuthHeaders(),
    })

    const data = await handleApiResponse<any>(response)

    return {
        enterprise: data.enterprise || null,
        stages: data.stages || [],
        tutors: data.tutors || []
    }
}

export async function addTutor(tutorData: Omit<Tutor, 'id' | 'enterprise_id' | 'created_at' | 'updated_at'>): Promise<Tutor> {
    const response = await fetch(`${API_URL}/api/enterprise/tutors`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify(tutorData)
    })

    return handleApiResponse<Tutor>(response)
}

export async function updateTutor(id: number, tutorData: Partial<Tutor>): Promise<Tutor> {
    const response = await fetch(`${API_URL}/api/enterprise/tutors/${id}`, {
        method: 'PUT',
        headers: getAuthHeaders(),
        body: JSON.stringify(tutorData)
    })

    return handleApiResponse<Tutor>(response)
}

export async function deleteTutor(id: number): Promise<void> {
    const response = await fetch(`${API_URL}/api/enterprise/tutors/${id}`, {
        method: 'DELETE',
        headers: getAuthHeaders(),
    })

    await handleApiResponse<void>(response)
}

// Nouvelles fonctions utilitaires
export async function getTutorById(id: number): Promise<Tutor> {
    const response = await fetch(`${API_URL}/api/enterprise/tutors/${id}`, {
        headers: getAuthHeaders(),
    })

    return handleApiResponse<Tutor>(response)
}

export async function getTutors(): Promise<Tutor[]> {
    const response = await fetch(`${API_URL}/api/enterprise/tutors`, {
        headers: getAuthHeaders(),
    })

    return handleApiResponse<Tutor[]>(response)
}
