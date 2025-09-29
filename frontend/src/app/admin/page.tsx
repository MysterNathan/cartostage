'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { authApi } from '@/lib/authApi'
import { getStages, addStage, updateStage, deleteStage } from '@/lib/stageApi'
import StageModal from '@/components/admin/StageModal'
import StatisticsModal from '@/components/admin/StatisticsModal'
import type {Stage, StageWithDetails} from '@/types/stage'
import FilieresManager from '@/components/admin/FilieresManager'
import StagesList from "@/components/admin/StageList";
import StatsCards from "@/components/admin/StatCard";

export default function AdminPage() {
  const router = useRouter()
  const [stages, setStages] = useState<Stage[]>([])
  const [editingStage, setEditingStage] = useState<Stage | null>(null)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [isStatsModalOpen, setIsStatsModalOpen] = useState(false)
  const [isNewStage, setIsNewStage] = useState(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [isFilieresOpen, setIsFilieresOpen] = useState(false)

  useEffect(() => {
    // Vérifier l'authentification au chargement
    if (!authApi.isAuthenticated()) {
      router.push('/login')
      return
    }

    loadStages()
  }, [router])

  const handleLogout = () => {
    authApi.logout()
    router.push('/login')
  }

  const loadStages = async () => {
    try {
      setLoading(true)
      setError('')

      const data: StageWithDetails[] = await getStages()
      setStages(data)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Impossible de charger les stages'
      setError(errorMessage)
      console.error('Erreur:', err)

      // Si l'erreur est liée à l'authentification, rediriger vers la page de login
      if (err instanceof Error && err.message.includes('401')) {
        authApi.logout()
        router.push('/login')
      }
    } finally {
      setLoading(false)
    }
  }

  // Affichage du loader pendant la vérification d'authentification/chargement
  if (loading) {
    return (
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-500 mx-auto mb-4"></div>
            <p className="text-gray-600">Chargement...</p>
          </div>
        </div>
    )
  }

// Dans votre page admin, garder seulement ces fonctions simplifiées :

  const editStage = (stage: Stage) => {
    setEditingStage(stage);
    setIsNewStage(false);
    setIsModalOpen(true);
  };

  const addNewStage = () => {
    // Stage vide pour la création
    setEditingStage(null); // Passer null, la modal gèrera l'initialisation
    setIsNewStage(true);
    setIsModalOpen(true);
  };

  const handleStageSuccess = (stage: Stage) => {
    if (isNewStage) {
      setStages(prev => [...prev, stage]);
    } else {
      setStages(prev => prev.map(s => s.id === stage.id ? stage : s));
    }

    // Fermer la modal
    setIsModalOpen(false);
    setEditingStage(null);
    setIsNewStage(false);
  };

  const handleModalClose = () => {
    setIsModalOpen(false);
    setEditingStage(null);
    setIsNewStage(false);
  };

  const totalStages = stages.length
  const totalCapacity = stages.reduce((sum, stage) => sum + stage.stage_offer.capacity_total, 0)
  const occupiedPlaces = stages.reduce((sum, stage) => sum + stage.stage_offer.capacity_filled, 0)
  const availablePlaces = totalCapacity - occupiedPlaces


return (
      <div className="min-h-screen bg-gray-100">
        {/* Header */}
        <div className="bg-white shadow-sm border-b">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between items-center h-16">
              <h1 className="text-xl font-semibold text-gray-900">Administration des stages</h1>
              <div className="flex items-center gap-4">
                <a
                    href="/"
                    className="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
                >
                  Retour à la carte
                </a>
                <button
                    onClick={() => setIsStatsModalOpen(true)}
                    className="flex items-center px-4 py-2 bg-purple-500 text-white rounded-lg hover:bg-purple-600 transition-colors"
                >
                  <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012-2z" />
                  </svg>
                  Statistiques
                </button>
                <button
                    onClick={() => setIsFilieresOpen(true)}
                    className="flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
                >
                  <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 7h18M3 12h18M3 17h18" />
                  </svg>
                  Gérer les filières
                </button>
                {/* Bouton pour ajouter un nouveau stage */}
                <button
                    onClick={addNewStage}
                    className="flex items-center px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                >
                  <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                  </svg>
                  Ajouter un stage
                </button>

                <button
                    onClick={handleLogout}
                    className="px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors"
                >
                  Déconnexion
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Contenu principal */}
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
                {error}
              </div>
          )}
          <StatsCards
              totalStages={totalStages}
              totalCapacity={totalCapacity}
              occupiedPlaces={occupiedPlaces}
              availablePlaces={availablePlaces}
          />
          <StagesList
              stages={stages}
              onEdit={editStage}
              //onDelete={handleDeleteStage}
              loading={loading}
          />
        </div>
        {/* Modal */}
        {isModalOpen && (
            <StageModal
                stage={editingStage}
                isOpen={isModalOpen}
                onClose={handleModalClose}
                onSuccess={handleStageSuccess}
                onDelete={(stageId) => setStages(stages.filter(s => s.id !== stageId))}
                isNew={isNewStage}
            />
        )}
        <FilieresManager
            isOpen={isFilieresOpen}
            onClose={() => setIsFilieresOpen(false)}
        />

        <StatisticsModal
            stages={stages}
            isOpen={isStatsModalOpen}
            onClose={() => setIsStatsModalOpen(false)}
        />
      </div>
  )
}
