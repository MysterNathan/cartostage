'use client'

import { useState, useEffect } from 'react'
import dynamic from 'next/dynamic'
import StageSelector from '@/components/StageSelector'
import type { Stage } from '@/types/stage'

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

  useEffect(() => {
    loadStages()
  }, [])

  const loadStages = async () => {
    try {
      setLoading(true)
      const response = await fetch('/api/stages')
      if (!response.ok) throw new Error('Erreur de chargement')
      
      const data = await response.json()
      const stagesData = data.stages || []
      
      setStages(stagesData)
      setFilteredStages(stagesData)
    } catch (error) {
      console.error('Erreur lors du chargement:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleStageClick = (stage: Stage) => {
    setSelectedStage(stage)
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
        <div className="w-80 bg-white border-r border-gray-200 overflow-hidden flex flex-col">
          <StageSelector 
            stages={stages}
            filteredStages={filteredStages}
            selectedStage={selectedStage}
            onStageSelect={handleStageClick}
            onFilterChange={setFilteredStages}
          />
        </div>
        
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
