// lib/stageApi.ts
import { Stage, StagesData } from "@/types/stage";

const API_URL = 'http://localhost:8080/api/stages'

export async function getStages(): Promise<StagesData> {
    const response = await fetch(API_URL, {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },
    })
    console.log(`response: ${response}`)
    if (!response.ok) {
        throw new Error(`Erreur lors de la récupération des stages: ${response.statusText}`)
    }

    const data: StagesData = await response.json()
    return data
}

export async function addStage(stage: Omit<Stage, 'id'>): Promise<Stage> {
    const response = await fetch(API_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(stage),
    })

    if (!response.ok) {
        throw new Error(`Erreur lors de l'ajout du stage: ${response.statusText}`)
    }

    const data: Stage = await response.json()
    return data
}

// lib/stageApi.ts
export async function updateStage(id: number, stage: Omit<Stage, 'id'>): Promise<Stage> {
    const response = await fetch(`${API_URL}/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
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
        headers: { 'Content-Type': 'application/json' },
    })

    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.error || `Erreur lors de la suppression du stage: ${response.statusText}`)
    }

    return await response.json()
}

