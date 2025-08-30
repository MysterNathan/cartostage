'use client'
import { useState } from 'react'
import StageFilters from './StageFilters'
import type { Stage } from '@/types/stage'

interface StageSelectorProps {
  stages: Stage[]
  filteredStages: Stage[]
  selectedStage?: Stage | null
  onStageSelect: (stage: Stage) => void
  onFilterChange: (filteredStages: Stage[]) => void
}

export default function StageSelector({ 
  stages, 
  filteredStages, 
  selectedStage, 
  onStageSelect,
  onFilterChange
}: StageSelectorProps) {
  const [isFiltersExpanded, setIsFiltersExpanded] = useState(true)

  return (
    <div className="h-full flex flex-col">
      {/* En-tête avec bouton de toggle */}
      <div className="shrink-0 bg-white border-b border-gray-200">
        <button
          onClick={() => setIsFiltersExpanded(!isFiltersExpanded)}
          className="w-full px-4 py-3 flex items-center justify-between hover:bg-gray-50 transition-colors"
        >
          <div className="flex items-center gap-2">
            <svg className="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.414A1 1 0 013 6.707V4z" />
            </svg>
            <span className="font-medium text-gray-900">Filtres</span>
            <span className="text-sm text-gray-500">
              ({filteredStages.length}/{stages.length})
            </span>
          </div>
          <svg 
            className={`w-5 h-5 text-gray-400 transition-transform ${
              isFiltersExpanded ? 'rotate-180' : ''
            }`} 
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
          </svg>
        </button>
      </div>

      {/* Section des filtres (rétractable) */}
      <div className={`shrink-0 overflow-hidden transition-all duration-300 ${
        isFiltersExpanded ? 'max-h-96' : 'max-h-0'
      }`}>
        <div className="border-b border-gray-200">
          <StageFilters 
            stages={stages}
            onFiltersChange={onFilterChange}
            totalStages={stages.length}
          />
        </div>
      </div>

      {/* Résumé compact des filtres actifs quand rétracté */}
      {!isFiltersExpanded && (
        <div className="shrink-0 px-4 py-2 bg-gray-50 border-b border-gray-200">
          <div className="flex items-center gap-4 text-sm text-gray-600">
            <span className="font-medium">
              {filteredStages.length} résultat{filteredStages.length > 1 ? 's' : ''}
            </span>
            {filteredStages.length !== stages.length && (
              <span className="text-blue-600">
                Filtres actifs
              </span>
            )}
          </div>
        </div>
      )}

      {/* Liste des stages */}
      <div className="flex-1 overflow-y-auto">
        <div className="p-4 space-y-3">
          {filteredStages.length === 0 ? (
            <div className="text-center py-8">
              <div className="text-gray-400 mb-2">
                <svg className="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9.172 16.172a4 4 0 015.656 0M9 12h6m-6-4h6m2 5.291A7.962 7.962 0 0112 15c-2.34 0-4.5-.935-6.052-2.473M15 21H9a2 2 0 01-2-2V5a2 2 0 012 2h6a2 2 0 012 2v14a2 2 0 01-2 2z" />
                </svg>
              </div>
              <p className="text-gray-500">Aucun stage trouvé</p>
              <p className="text-sm text-gray-400">Essayez de modifier vos critères de recherche</p>
            </div>
          ) : (
            filteredStages.map((stage) => (
              <div
                key={stage.id}
                onClick={() => onStageSelect(stage)}
                className={`p-4 rounded-lg border cursor-pointer transition-all hover:shadow-md ${
                  selectedStage?.id === stage.id
                    ? 'border-blue-500 bg-blue-50 shadow-sm'
                    : 'border-gray-200 bg-white hover:border-gray-300'
                }`}
              >
                {/* En-tête avec filière et disponibilité */}
                <div className="flex justify-between items-start mb-2">
                  <span className={`px-2 py-1 text-xs font-medium rounded-full ${
                    stage.filiere === 'CCST' 
                      ? 'bg-blue-100 text-blue-800' 
                      : 'bg-green-100 text-green-800'
                  }`}>
                    {stage.filiere}
                  </span>
                  <div className="text-right">
                    <div className={`text-sm font-medium ${
                      stage.placesDisponibles > 0 ? 'text-green-600' : 'text-red-600'
                    }`}>
                      {stage.placesDisponibles > 0 
                        ? `${stage.placesDisponibles} place${stage.placesDisponibles > 1 ? 's' : ''}`
                        : 'Complet'
                      }
                    </div>
                    <div className="text-xs text-gray-500">
                      {stage.capacity_filled}/{stage.capacity_total}
                    </div>
                  </div>
                </div>

                {/* Informations principales */}
                <div className="mb-2">
                  <h4 className="font-semibold text-gray-900 text-sm leading-tight mb-1">
                    {stage.poste}
                  </h4>
                  <p className="text-sm text-gray-700 font-medium">
                    {stage.entreprise}
                  </p>
                  {stage.sector && (
                    <p className="text-xs text-gray-600">
                      {stage.sector}
                    </p>
                  )}
                </div>

                {/* Localisation */}
                <div className="flex items-center text-xs text-gray-500 mb-2">
                  <svg className="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  {stage.commune}
                </div>

                {/* Périodes si disponibles */}
                {stage.period && (
                  <div className="flex items-center text-xs text-gray-500">
                    <svg className="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 002 2z" />
                    </svg>
                    {stage.period}
                  </div>
                )}

                {/* Indicateur de sélection */}
                {selectedStage?.id === stage.id && (
                  <div className="mt-2 pt-2 border-t border-blue-200">
                    <div className="flex items-center text-xs text-blue-600">
                      <svg className="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                      </svg>
                      Stage sélectionné
                    </div>
                  </div>
                )}
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  )
}
