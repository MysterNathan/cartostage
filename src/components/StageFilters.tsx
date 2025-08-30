'use client'
import { useState } from 'react'
import type { Stage } from '@/types/stage'

interface StageFiltersProps {
  stages: Stage[]
  onFiltersChange: (filteredStages: Stage[]) => void
  totalStages: number
}

export default function StageFilters({ 
  stages, 
  onFiltersChange, 
  totalStages 
}: StageFiltersProps) {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedFiliere, setSelectedFiliere] = useState('')
  const [selectedCommune, setSelectedCommune] = useState('')
  const [minPlaces, setMinPlaces] = useState(0)
  const [onlyAvailable, setOnlyAvailable] = useState(false)

  // Extraire les options uniques
  const filieres = Array.from(new Set(stages.map(s => s.filiere))).sort()
  const communes = Array.from(new Set(stages.map(s => s.commune))).sort()

  // Fonction de filtrage
  const applyFilters = (
    query: string,
    filiere: string,
    commune: string,
    places: number,
    available: boolean
  ) => {
    let filtered = stages.filter(stage => {
      const matchesQuery = !query || 
        [stage.poste, stage.entreprise, stage.sector, stage.commune, stage.adresse]
          .join(' ').toLowerCase().includes(query.toLowerCase())
      
      const matchesFiliere = !filiere || stage.filiere === filiere
      const matchesCommune = !commune || stage.commune === commune
      const matchesPlaces = stage.placesDisponibles >= places
      const matchesAvailable = !available || stage.placesDisponibles > 0

      return matchesQuery && matchesFiliere && matchesCommune && matchesPlaces && matchesAvailable
    })

    onFiltersChange(filtered)
  }

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const query = e.target.value
    setSearchQuery(query)
    applyFilters(query, selectedFiliere, selectedCommune, minPlaces, onlyAvailable)
  }

  const handleFiliereChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const filiere = e.target.value
    setSelectedFiliere(filiere)
    applyFilters(searchQuery, filiere, selectedCommune, minPlaces, onlyAvailable)
  }

  const handleCommuneChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const commune = e.target.value
    setSelectedCommune(commune)
    applyFilters(searchQuery, selectedFiliere, commune, minPlaces, onlyAvailable)
  }

  const handlePlacesChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const places = parseInt(e.target.value) || 0
    setMinPlaces(places)
    applyFilters(searchQuery, selectedFiliere, selectedCommune, places, onlyAvailable)
  }

  const handleAvailableChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const available = e.target.checked
    setOnlyAvailable(available)
    applyFilters(searchQuery, selectedFiliere, selectedCommune, minPlaces, available)
  }

  const clearFilters = () => {
    setSearchQuery('')
    setSelectedFiliere('')
    setSelectedCommune('')
    setMinPlaces(0)
    setOnlyAvailable(false)
    onFiltersChange(stages)
  }

  // Calculer les statistiques
  const filteredStages = stages.filter(stage => {
    const matchesQuery = !searchQuery || 
      [stage.poste, stage.entreprise, stage.sector, stage.commune, stage.adresse]
        .join(' ').toLowerCase().includes(searchQuery.toLowerCase())
    
    const matchesFiliere = !selectedFiliere || stage.filiere === selectedFiliere
    const matchesCommune = !selectedCommune || stage.commune === selectedCommune
    const matchesPlaces = stage.placesDisponibles >= minPlaces
    const matchesAvailable = !onlyAvailable || stage.placesDisponibles > 0

    return matchesQuery && matchesFiliere && matchesCommune && matchesPlaces && matchesAvailable
  })

  const totalCapacity = filteredStages.reduce((sum, stage) => sum + stage.capacity_total, 0)
  const totalFilled = filteredStages.reduce((sum, stage) => sum + stage.capacity_filled, 0)
  const totalAvailable = filteredStages.reduce((sum, stage) => sum + stage.placesDisponibles, 0)

  return (
    <div className="bg-white p-4 rounded-lg shadow-sm border">
      <h3 className="font-semibold mb-3">Filtrer les stages</h3>
      
      <div className="space-y-3">
        {/* Recherche */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Rechercher
          </label>
          <input
            type="text"
            value={searchQuery}
            onChange={handleSearchChange}
            placeholder="Poste, entreprise, secteur..."
            className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        {/* Filière */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Filière
          </label>
          <select
            value={selectedFiliere}
            onChange={handleFiliereChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">Toutes filières</option>
            {filieres.map(filiere => (
              <option key={filiere} value={filiere}>
                {filiere}
              </option>
            ))}
          </select>
        </div>

        {/* Commune */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Commune
          </label>
          <select
            value={selectedCommune}
            onChange={handleCommuneChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">Toutes communes</option>
            {communes.map(commune => (
              <option key={commune} value={commune}>
                {commune}
              </option>
            ))}
          </select>
        </div>

        {/* Places minimum */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Places libres minimum: {minPlaces}
          </label>
          <input
            type="range"
            min="0"
            max="5"
            value={minPlaces}
            onChange={handlePlacesChange}
            className="w-full"
          />
        </div>

        {/* Uniquement disponibles */}
        <div className="flex items-center">
          <input
            type="checkbox"
            id="onlyAvailable"
            checked={onlyAvailable}
            onChange={handleAvailableChange}
            className="mr-2"
          />
          <label htmlFor="onlyAvailable" className="text-sm text-gray-700">
            Uniquement stages avec places libres
          </label>
        </div>

        <button
          onClick={clearFilters}
          className="w-full px-3 py-2 bg-gray-100 text-gray-700 rounded-md text-sm hover:bg-gray-200 transition-colors"
        >
          Effacer les filtres
        </button>
      </div>

      {/* Statistiques */}
      <div className="mt-4 pt-4 border-t border-gray-200">
        <div className="text-sm text-gray-600 space-y-1">
          <div className="flex justify-between">
            <span>Résultats:</span>
            <span className="font-medium">{filteredStages.length} / {totalStages}</span>
          </div>
          <div className="flex justify-between">
            <span>Capacité totale:</span>
            <span className="font-medium">{totalCapacity}</span>
          </div>
          <div className="flex justify-between">
            <span>Places occupées:</span>
            <span className="font-medium">{totalFilled}</span>
          </div>
          <div className="flex justify-between">
            <span>Places disponibles:</span>
            <span className="font-medium text-green-600">{totalAvailable}</span>
          </div>
          {totalCapacity > 0 && (
            <div className="flex justify-between">
              <span>Taux d'occupation:</span>
              <span className="font-medium">{Math.round((totalFilled / totalCapacity) * 100)}%</span>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
