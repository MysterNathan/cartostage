// lib/enterpriseApi.ts
import type { EnterpriseData } from '@/types/enterprise'
import type { Stage } from '@/types/stage'
import type { Tutor } from '@/types/tutor'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export async function getMyEnterpriseData(): Promise<EnterpriseData> {
    const token = localStorage.getItem('authToken')
    if (!token) {
        throw new Error('Token d\'authentification manquant')
    }

    const response = await fetch(`${API_URL}/api/enterprises/me`, {
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
        },
    })

    if (!response.ok) {
        if (response.status === 401) {
            throw new Error('Non autorisé - Token invalide')
        }
        throw new Error(`Erreur API: ${response.status}`)
    }

    const data = await response.json()

    return {
        enterprise: data.enterprise || null,
        stages: data.stages || [],
        tutors: data.tutors || []
    }
}

export async function addTutor(tutorData: Omit<Tutor, 'id' | 'enterprise_id' | 'created_at' | 'updated_at'>): Promise<Tutor> {
    const token = localStorage.getItem('authToken')
    if (!token) {
        throw new Error('Token d\'authentification manquant')
    }

    const response = await fetch(`${API_URL}/api/enterprise/tutors`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(tutorData)
    })

    if (!response.ok) {
        throw new Error(`Erreur lors de l'ajout du tuteur: ${response.status}`)
    }

    return response.json()
}

export async function updateTutor(id: number, tutorData: Partial<Tutor>): Promise<Tutor> {
    const token = localStorage.getItem('authToken')
    if (!token) {
        throw new Error('Token d\'authentification manquant')
    }

    const response = await fetch(`${API_URL}/api/enterprise/tutors/${id}`, {
        method: 'PUT',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(tutorData)
    })

    if (!response.ok) {
        throw new Error(`Erreur lors de la modification du tuteur: ${response.status}`)
    }

    return response.json()
}

export async function deleteTutor(id: number): Promise<void> {
    const token = localStorage.getItem('authToken')
    if (!token) {
        throw new Error('Token d\'authentification manquant')
    }

    const response = await fetch(`${API_URL}/api/enterprise/tutors/${id}`, {
        method: 'DELETE',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
        },
    })

    if (!response.ok) {
        throw new Error(`Erreur lors de la suppression du tuteur: ${response.status}`)
    }
}
