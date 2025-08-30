'use client'
import { useState, useMemo } from 'react'
import type { Stage } from '@/types/stage'

export function useStageFilters(stages: Stage[]) {
  const [searchQuery, setSearchQuery] = useState('')
  const [minPlaces, setMinPlaces] = useState(0)

  const filteredStages = useMemo(() => {
    return stages.filter(stage => {
      // Filtre par recherche textuelle
      const matchesSearch = !searchQuery || 
        stage.poste.toLowerCase().includes(searchQuery.toLowerCase()) ||
        stage.entreprise.toLowerCase().includes(searchQuery.toLowerCase()) ||
        stage.adresse.toLowerCase().includes(searchQuery.toLowerCase())

      // Filtre par places minimum
      const matchesPlaces = stage.placesDisponibles >= minPlaces

      return matchesSearch && matchesPlaces
    })
  }, [stages, searchQuery, minPlaces])

  return {
    filteredStages,
    searchQuery,
    minPlaces,
    setSearchQuery,
    setMinPlaces,
    totalStages: stages.length,
    filteredCount: filteredStages.length
  }
}
