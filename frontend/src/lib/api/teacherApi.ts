const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

import type { Tutor } from '@/types/tutor'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

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

    // Vérifier s'il y a du contenu à parser
    const contentType = response.headers.get('content-type')
    const contentLength = response.headers.get('content-length')

    // Si pas de contenu ou contenu vide, retourner undefined
    if (response.status === 204 || contentLength === '0' || !contentType?.includes('application/json')) {
        return undefined as T
    }

    try {
        return await response.json()
    } catch (error) {
        // Si erreur de parsing, probablement une réponse vide
        console.warn('Réponse vide ou invalide, retour de undefined')
        return undefined as T
    }
}

export async function getTutors(): Promise<Tutor[]> {
    const response = await fetch(`${API_URL}/api/tutors`, {
        headers: getAuthHeaders(),
    })

    return handleApiResponse<Tutor[]>(response)
}