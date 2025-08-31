import { useState, useEffect } from 'react'
import type { Stage, StagesData } from '@/types/stage'

export function useStagesAdmin() {
  const [stages, setStages] = useState<Stage[]>([])
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadStages()
  }, [])

  const loadStages = async () => {
    try {
      setLoading(true)
      const response = await fetch('/api/stages')
      
      if (!response.ok) {
        throw new Error('Erreur de chargement')
      }
      
      const data: StagesData = await response.json()
      setStages(data.stages)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur inconnue')
    } finally {
      setLoading(false)
    }
  }

  const saveStages = async (newStages: Stage[]) => {
    try {
      setSaving(true)
      setError(null)
      
      const response = await fetch('/api/stages', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ stages: newStages }),
      })

      if (!response.ok) {
        throw new Error('Erreur de sauvegarde')
      }

      setStages(newStages)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur de sauvegarde')
      throw err
    } finally {
      setSaving(false)
    }
  }

  return {
    stages,
    loading,
    saving,
    error,
    saveStages,
    reload: loadStages,
  }
}
