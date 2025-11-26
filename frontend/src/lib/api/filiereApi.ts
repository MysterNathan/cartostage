import type { Filiere, FilieresResponse, CreateFiliereData, UpdateFiliereData } from '@/types/filiere'
import authApi from "@/lib/api/authApi";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;
const API_URL = `${API_BASE_URL}/api/filieres`;

// Utilitaire pour récupérer les headers avec authentification
const getAuthHeaders = (): HeadersInit => {
    const token = authApi.getToken();
    return {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
    };
};

export async function getFilieres(): Promise<FilieresResponse> {
    const response = await fetch(API_URL, {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
    })

    if (!response.ok) {
        throw new Error(`Erreur lors de la récupération des filières: ${response.statusText}`)
    }

    const data = await response.json()
    return data
}

export async function getFiliere(id: number): Promise<Filiere> {
    const response = await fetch(`${API_URL}/${id}`, {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
    })

    if (!response.ok) {
        throw new Error(`Erreur lors de la récupération de la filière: ${response.statusText}`)
    }

    return await response.json()
}

export async function addFiliere(filiere: CreateFiliereData): Promise<Filiere> {
    const response = await fetch(API_URL, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify(filiere),
    })

    if (!response.ok) {
        throw new Error(`Erreur lors de la création de la filière: ${response.statusText}`)
    }

    return await response.json()
}

export async function updateFiliere(id: number, filiere: UpdateFiliereData): Promise<Filiere> {
    const response = await fetch(`${API_URL}/${id}`, {
        method: 'PUT',
        headers: getAuthHeaders(),
        body: JSON.stringify(filiere),
    })

    if (!response.ok) {
        throw new Error(`Erreur lors de la modification de la filière: ${response.statusText}`)
    }

    return await response.json()
}

export async function deleteFiliere(id: number): Promise<void> {
    const response = await fetch(`${API_URL}/${id}`, {
        method: 'DELETE',
        headers: getAuthHeaders(),
    })

    if (!response.ok) {
        throw new Error(`Erreur lors de la suppression de la filière: ${response.statusText}`)
    }
}

// Fonction utilitaire pour obtenir une filière par son code
export async function getFiliereByCode(code: string): Promise<Filiere | null> {
    try {
        const { filieres } = await getFilieres()
        return filieres.find(f => f.code === code) || null
    } catch (error) {
        console.error('Erreur lors de la recherche de filière par code:', error)
        return null
    }
}

// Fonction utilitaire pour obtenir les couleurs des filières sous forme de map
export async function getFilieresColorMap(): Promise<Record<string, string>> {
    try {
        const { filieres } = await getFilieres()
        return filieres.reduce((acc, filiere) => {
            acc[filiere.code] = filiere.color
            return acc
        }, {} as Record<string, string>)
    } catch (error) {
        console.error('Erreur lors de la récupération des couleurs:', error)
        return {}
    }
}
