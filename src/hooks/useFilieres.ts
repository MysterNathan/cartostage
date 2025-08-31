'use client'
import { useCallback, useEffect, useState } from 'react'
import type { Filiere } from '@/types/filiere'


interface CreateInput { code: string; label: string; color?: string }
interface UpdateInput { id: number; code?: string; label?: string; color?: string }


export function useFilieres() {
    const [filieres, setFilieres] = useState<Filiere[]>([])
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)


    const refresh = useCallback(async () => {
        try {
            setLoading(true)
            setError(null)
            const r = await fetch('/api/filieres', { cache: 'no-store' })
            if (!r.ok) throw new Error('GET /api/filieres failed')
            const d = await r.json()
            setFilieres((d?.filieres ?? []).sort((a: Filiere, b: Filiere) => a.label.localeCompare(b.label)))
        } catch (e: any) {
            setError(e?.message ?? 'Erreur chargement filières')
        } finally {
            setLoading(false)
        }
    }, [])


    const create = useCallback(async (payload: CreateInput) => {
        const r = await fetch('/api/filieres', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        })
        if (!r.ok) throw new Error('POST /api/filieres failed')
        const d = await r.json()
        const f = d?.filiere as Filiere
        setFilieres(prev => [...prev.filter(x => x.id !== f.id), f].sort((a, b) => a.label.localeCompare(b.label)))
        return f
    }, [])


    const update = useCallback(async (payload: UpdateInput) => {
        const r = await fetch('/api/filieres', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        })
        if (!r.ok) throw new Error('PUT /api/filieres failed')
        const d = await r.json()
        const f = d?.filiere as Filiere
        setFilieres(prev => prev.map(x => (x.id === f.id ? f : x)).sort((a, b) => a.label.localeCompare(b.label)))
        return f
    }, [])


    const remove = useCallback(async (id: number) => {
        const r = await fetch('/api/filieres', {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id })
        })
        if (!r.ok) throw new Error('DELETE /api/filieres failed')
        setFilieres(prev => prev.filter(x => x.id !== id))
        return true
    }, [])


    useEffect(() => { refresh() }, [refresh])


    return { filieres, loading, error, refresh, create, update, remove, setFilieres }
}