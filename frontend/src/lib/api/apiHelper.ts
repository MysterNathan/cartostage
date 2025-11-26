// lib/api/apiHelpers.ts
const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

/**
 * Récupère les headers d'authentification
 */
export function getAuthHeaders(): HeadersInit {
    const token = localStorage.getItem('authToken')
    if (!token) {
        throw new Error('Token d\'authentification manquant')
    }

    return {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
    }
}

/**
 * Gère les réponses d'API avec traitement des erreurs
 */
export async function handleApiResponse<T>(response: Response): Promise<T> {
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

    const contentType = response.headers.get('content-type')
    const contentLength = response.headers.get('content-length')

    if (response.status === 204 || contentLength === '0' || !contentType?.includes('application/json')) {
        return undefined as T
    }

    try {
        return await response.json()
    } catch (error) {
        console.warn('Réponse vide ou invalide, retour de undefined')
        return undefined as T
    }
}

/**
 * Effectue une requête API avec gestion des erreurs
 */
export async function fetchApi<T>(
    endpoint: string,
    options?: RequestInit
): Promise<T> {
    const response = await fetch(`${API_URL}${endpoint}`, {
        ...options,
        headers: {
            ...getAuthHeaders(),
            ...options?.headers,
        },
    })

    return handleApiResponse<T>(response)
}

export { API_URL }
