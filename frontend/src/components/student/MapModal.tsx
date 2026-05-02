// components/carte/StageMapView.tsx
'use client'

import { useState, useEffect } from 'react'
import dynamic from 'next/dynamic'
import StageSelector from '@/components/carte/StageSelector'
import { getStages } from '@/lib/api/stageApi'
import type { Stage, StageWithDetails } from '@/types/stage'

// Import dynamique pour éviter l'erreur SSR avec Leaflet
const StageMap = dynamic(() => import('@/components/carte/StageMap'), {
    ssr: false,
    loading: () => (
        <div className="h-96 bg-gray-100 rounded-lg flex items-center justify-center">
            <div className="text-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-2"></div>
                <p className="text-gray-700">Chargement de la carte...</p>
            </div>
        </div>
    ),
})

interface StageMapViewProps {
    title?: string
    subtitle?: string
    showAdminLink?: boolean
    stages?: Stage[]
}

export default function StageMapView({
                                         stages,
                                     }: StageMapViewProps) {
    const [filteredStages, setFilteredStages] = useState<Stage[]>(stages || [])
    const [selectedStage, setSelectedStage] = useState<Stage | null>(null)

    const handleStageClick = (stage: Stage) => {
        setSelectedStage(stage)
    }

    // === Affichage normal ===
    return (
        <>
        {/* Contenu */}
    <div className="flex h-[calc(100vh-80px)]">

        {/* Carte */}
        <div className="flex-1 border-l border-gray-200">
            <StageMap
                stages={filteredStages}
                selectedStage={selectedStage}
                onStageClick={handleStageClick}
            />
        </div>
    </div>
        </>
    )
}
