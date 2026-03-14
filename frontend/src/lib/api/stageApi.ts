// lib/stageApi.ts
import { Stage, StagesData } from "@/types/stage";
import { authApi } from "@/lib/api/authApi";
import {FormResponse, FormResponses, FormSection} from "@/types/form";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;
const API_URL_PUBLIC = `${API_BASE_URL}/api/stages/public`;
const API_URL = `${API_BASE_URL}/api/stages`;

// Utilitaire pour récupérer les headers avec authentification
const getAuthHeaders = (): HeadersInit => {
    const token = authApi.getToken();
    return {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
    };
};

export async function getStages(): Promise<StagesData> {
    const response = await fetch(API_URL_PUBLIC, {
        method: 'GET',
        headers: getAuthHeaders(),
    })

    if (!response.ok) {
        throw new Error(`Erreur lors de la récupération des stages: ${response.statusText}`)
    }

    const data: StagesData = await response.json()
    return data
}

export async function addStage(stage: Omit<Stage, 'id'>): Promise<Stage> {
    const response = await fetch(API_URL, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify(stage),
    })

    if (!response.ok) {
        throw new Error(`Erreur lors de l'ajout du stage: ${response.statusText}`)
    }

    const data: Stage = await response.json()
    return data
}

export async function updateStage(id: number, stage: Omit<Stage, 'id'>): Promise<Stage> {
    const response = await fetch(`${API_URL}/${id}`, {
        method: 'PUT',
        headers: getAuthHeaders(),
        body: JSON.stringify(stage),
    })

    if (!response.ok) {
        throw new Error(`Erreur lors de la modification du stage: ${response.statusText}`)
    }

    return await response.json()
}

export async function deleteStage(id: number): Promise<{ success: boolean; message: string }> {
    const response = await fetch(`${API_URL}/${id}`, {
        method: 'DELETE',
        headers: getAuthHeaders(),
    })

    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.error || `Erreur lors de la suppression du stage: ${response.statusText}`)
    }

    return await response.json()
}

export async function getForms(): Promise<FormResponses> {
    const response = await fetch(`${API_URL}/form`, {
        method: 'GET',
        headers: getAuthHeaders(),
    })
    return await response.json()
}
