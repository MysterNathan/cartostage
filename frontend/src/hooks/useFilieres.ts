'use client'
import { useState, useEffect, useCallback } from 'react'
import { getFilieres, addFiliere, updateFiliere, deleteFiliere } from '@/lib/filiereApi'
import type { Filiere } from '@/types/filiere'

export interface UseFilieres {
    // État
    filieres: Filiere[]
    loading: boolean
    error: string | null

    // Actions
    loadFilieres: () => Promise<void>
    create: (data: { code: string; label: string; color: string }) => Promise<void>
    update: (data: { id: number; code?: string; label?: string; color?: string }) => Promise<void>
    remove: (id: number) => Promise<void>
}

export function useFilieres(): UseFilieres {
    const [filieres, setFilieres] = useState<Filiere[]>([])
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)

    // Utiliser useCallback pour éviter la re-création de la fonction à chaque rendu
    const loadFilieres = useCallback(async () => {
        try {
            setLoading(true)
            setError(null)
            const response = await getFilieres()
            setFilieres(response.filieres ?? [])
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Erreur lors du chargement')
        } finally {
            setLoading(false)
        }
    }, [])

    // Créer une filière
    const create = useCallback(async (data: { code: string; label: string; color: string }) => {
        try {
            setLoading(true)
            setError(null)
            const newFiliere = await addFiliere(data)
            setFilieres(prev => [...prev, newFiliere])
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Erreur lors de la création')
            throw err
        } finally {
            setLoading(false)
        }
    }, [])

    // Mettre à jour une filière
    const update = useCallback(async (data: { id: number; code?: string; label?: string; color?: string }) => {
        try {
            setLoading(true)
            setError(null)
            const updatedFiliere = await updateFiliere(data.id, data)
            setFilieres(prev => prev.map(f => f.id === data.id ? updatedFiliere : f))
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Erreur lors de la mise à jour')
            throw err
        } finally {
            setLoading(false)
        }
    }, [])

    // Supprimer une filière
    const remove = useCallback(async (id: number) => {
        try {
            setLoading(true)
            setError(null)
            await deleteFiliere(id)
            setFilieres(prev => prev.filter(f => f.id !== id))
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Erreur lors de la suppression')
            throw err
        } finally {
            setLoading(false)
        }
    }, [])

    return {
        filieres,
        loading,
        error,
        loadFilieres,
        create,
        update,
        remove
    }
}
