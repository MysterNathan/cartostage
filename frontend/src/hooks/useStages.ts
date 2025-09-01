import { useState, useEffect } from 'react'
import type { Stage, StagesData } from '@/types/stage'

export function useStages() {
  const [stages, setStages] = useState<Stage[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function loadStages() {
      try {
        setLoading(true)
        const response = await fetch('/api/stages')
        
        if (!response.ok) {
          throw new Error(`Erreur ${response.status}: ${response.statusText}`)
        }
        
        const data: StagesData = await response.json()
        setStages(data.stages)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Erreur inconnue')
      } finally {
        setLoading(false)
      }
    }

    loadStages()
  }, [])

  return { stages, loading, error }
}
