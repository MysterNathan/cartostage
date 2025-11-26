'use client'

import { useState, useEffect } from 'react'
import dynamic from 'next/dynamic'
import StageSelector from '@/components/StageSelector'
import { getStages } from '@/lib/api/stageApi'
import type {Stage, StageWithDetails} from '@/types/stage'

// Import dynamique pour éviter les erreurs SSR avec Leaflet
const StageMap = dynamic(() => import('@/components/StageMap'), {
  ssr: false,
  loading: () => (
      <div className="h-96 bg-gray-100 rounded-lg flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-2"></div>
          <p className="text-gray-600">Chargement de la carte...</p>
        </div>
      </div>
  )
})

export default function HomePage() {
  const [stages, setStages] = useState<Stage[]>([])
  const [filteredStages, setFilteredStages] = useState<Stage[]>([])
  const [selectedStage, setSelectedStage] = useState<Stage | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadStages()
  }, [])

  const loadStages = async () => {
    try {
      setLoading(true)
      setError(null)

      const data: StageWithDetails[] = await getStages()

      setStages(data)
      setFilteredStages(data)
    } catch (error) {
      console.error('Erreur lors du chargement:', error)
      setError(error instanceof Error ? error.message : 'Erreur inconnue')
    } finally {
      setLoading(false)
    }
  }

  const handleStageClick = (stage: Stage) => {
    setSelectedStage(stage)
  }

  const handleRetry = () => {
    loadStages()
  }

  if (loading) {
    return (
        <div className="min-h-screen bg-gray-50 flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
            <p className="text-gray-600">Chargement des stages...</p>
          </div>
        </div>
    )
  }

  if (error) {
    return (
        <div className="min-h-screen bg-gray-50 flex items-center justify-center">
          <div className="text-center max-w-md mx-auto p-6">
            <div className="text-red-500 mb-4">
              <svg className="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
              </svg>
            </div>
            <h2 className="text-xl font-semibold text-gray-900 mb-2">Erreur de chargement</h2>
            <p className="text-gray-600 mb-4">{error}</p>
            <button
                onClick={handleRetry}
                className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
            >
              Réessayer
            </button>
          </div>
        </div>
    )
  }

  return (
      <div className="min-h-screen bg-gray-50">
        <header className="bg-white shadow-sm border-b border-gray-200">
          <div className="px-6 py-4">
            <div className="flex justify-between items-center">
              <div>
                <h1 className="text-2xl font-bold text-gray-900">CartoStages Réunion Sud</h1>
                <p className="text-gray-600">CCST & ASSP - Lycées du Sud de La Réunion</p>
              </div>
              <div className="flex gap-3">
                <div className="text-sm text-gray-600">
                  <span className="font-medium">{stages.length}</span> stages disponibles
                </div>
                <button
                    onClick={handleRetry}
                    className="px-3 py-2 text-sm text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
                    title="Actualiser les données"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                </button>
                <a
                    href="/admin"
                    className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors text-sm"
                >
                  Administration
                </a>
              </div>
            </div>
          </div>
        </header>

        <div className="flex h-[calc(100vh-80px)]">
          {/* Panneau latéral avec rétractation */}
          <StageSelector
              stages={stages}
              filteredStages={filteredStages}
              selectedStage={selectedStage}
              onStageSelect={handleStageClick}
              onFilterChange={setFilteredStages}
          />

          {/* Zone de la carte */}
          <div className="flex-1">
            <StageMap
                stages={filteredStages}
                selectedStage={selectedStage}
                onStageClick={handleStageClick}
            />
          </div>
        </div>
      </div>
  )
}
