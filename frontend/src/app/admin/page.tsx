'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { authApi } from '@/lib/authApi'
import { getStages, addStage, updateStage, deleteStage } from '@/lib/stageApi'
import StageModal from '@/components/admin/StageModal'
import StatisticsModal from '@/components/admin/StatisticsModal'
import type { Stage } from '@/types/stage'
import FilieresManager from '@/components/admin/FilieresManager'

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

      const data = await getStages()
      setStages(data.stages || [])
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

  const editStage = (stage: Stage) => {
    setEditingStage(stage)
    setIsNewStage(false)
    setIsModalOpen(true)
  }

  const addNewStage = () => {
    const newStage: Stage = {
      id: Date.now(),
      entreprise: '',
      poste: '',
      adresse: '',
      commune: '',
      sector: '',
      filiere: 'CCST',
      lat: -20.8789,
      lng: 55.4481,
      capacity_total: 1,
      capacity_filled: 0,
      placesDisponibles: 1,
      period: ''
    }
    setEditingStage(newStage)
    setIsNewStage(true)
    setIsModalOpen(true)
  }

  const handleDeleteStage = async (stageId: number) => {
    const stage = stages.find(stage => stage.id === stageId)

    if (!confirm(`Êtes-vous sûr de vouloir supprimer le stage "${stage?.poste || 'inconnu'}" ?`)) {
      return
    }

    try {
      await deleteStage(stageId)
      setStages(stages.filter(s => s.id !== stageId))
      console.log(`Stage "${stage?.poste || 'inconnu'}" supprimé avec succès`)
    } catch (error) {
      console.error('Erreur suppression:', error)

      // Vérifier si c'est une erreur d'authentification
      if (error instanceof Error && error.message.includes('401')) {
        authApi.logout()
        router.push('/login')
        return
      }

      alert(error instanceof Error ? error.message : 'Erreur lors de la suppression')
    }
  }

  const handleStageSuccess = (stage: Stage) => {
    try {
      setError('')

      if (isNewStage) {
        setStages(prev => [...prev, stage])
        console.log('Nouveau stage ajouté à la liste')
      } else {
        setStages(prev => prev.map(s => s.id === stage.id ? stage : s))
        console.log('Stage mis à jour dans la liste')
      }

      setIsModalOpen(false)
      setEditingStage(null)

    } catch (error) {
      console.error('Erreur mise à jour état:', error)
      setError('Erreur lors de la mise à jour de l\'affichage')
    }
  }

  const totalStages = stages.length
  const totalCapacity = stages.reduce((sum, stage) => sum + stage.capacity_total, 0)
  const occupiedPlaces = stages.reduce((sum, stage) => sum + stage.capacity_filled, 0)
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

          {/* Statistiques rapides */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            <div className="bg-white p-6 rounded-lg shadow-sm border">
              <div className="flex items-center">
                <div className="p-2 bg-blue-100 rounded-lg">
                  <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                  </svg>
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-500">Total stages</p>
                  <p className="text-2xl font-semibold text-gray-900">{totalStages}</p>
                </div>
              </div>
            </div>

            <div className="bg-white p-6 rounded-lg shadow-sm border">
              <div className="flex items-center">
                <div className="p-2 bg-green-100 rounded-lg">
                  <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                  </svg>
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-500">Capacité totale</p>
                  <p className="text-2xl font-semibold text-gray-900">{totalCapacity}</p>
                </div>
              </div>
            </div>

            <div className="bg-white p-6 rounded-lg shadow-sm border">
              <div className="flex items-center">
                <div className="p-2 bg-orange-100 rounded-lg">
                  <svg className="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
                  </svg>
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-500">Places occupées</p>
                  <p className="text-2xl font-semibold text-gray-900">{occupiedPlaces}</p>
                </div>
              </div>
            </div>

            <div className="bg-white p-6 rounded-lg shadow-sm border">
              <div className="flex items-center">
                <div className="p-2 bg-purple-100 rounded-lg">
                  <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192L5.636 18.364M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-5 0a4 4 0 11-8 0 4 4 0 018 0z" />
                  </svg>
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-500">Places disponibles</p>
                  <p className="text-2xl font-semibold text-gray-900">{availablePlaces}</p>
                </div>
              </div>
            </div>
          </div>

          {/* Liste des stages */}
          <div className="bg-white shadow-sm rounded-lg border">
            <div className="px-6 py-4 border-b border-gray-200">
              <h3 className="text-lg font-medium text-gray-900">
                Liste des stages ({stages.length})
              </h3>
            </div>

            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Entreprise
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Poste
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Commune
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Filière
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Parcours
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Niveau
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Famille Métiers
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Capacité
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                {stages.map((stage) => (
                    <tr key={stage.id} className="hover:bg-gray-50">
                      <td className="px-4 py-4 whitespace-nowrap">
                        <div className="text-sm font-medium text-gray-900">{stage.entreprise}</div>
                        <div className="text-xs text-gray-500">{stage.sector || 'Secteur non renseigné'}</div>
                      </td>
                      <td className="px-4 py-4">
                        <div className="text-sm text-gray-900 max-w-xs truncate" title={stage.poste}>
                          {stage.poste}
                        </div>
                      </td>
                      <td className="px-4 py-4 whitespace-nowrap">
                        <div className="text-sm text-gray-900">{stage.commune || 'Non renseigné'}</div>
                      </td>
                      <td className="px-4 py-4 whitespace-nowrap">
              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                  stage.filiere === 'CCST' ? 'bg-blue-100 text-blue-800' :
                      stage.filiere === 'SN' ? 'bg-green-100 text-green-800' :
                          stage.filiere === 'MELEC' ? 'bg-yellow-100 text-yellow-800' :
                              stage.filiere === 'TMA' ? 'bg-red-100 text-red-800' :
                                  'bg-purple-100 text-purple-800'
              }`}>
                {stage.filiere}
              </span>
                      </td>
                      <td className="px-4 py-4 whitespace-nowrap">
              <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${
                  stage.parcours === 'scolaire' ? 'bg-blue-50 text-blue-700 border border-blue-200' :
                      stage.parcours === 'apprentissage' ? 'bg-orange-50 text-orange-700 border border-orange-200' :
                          stage.parcours === 'mixte' ? 'bg-purple-50 text-purple-700 border border-purple-200' :
                              'bg-gray-50 text-gray-700 border border-gray-200'
              }`}>
                {stage.parcours || 'Non défini'}
              </span>
                      </td>
                      <td className="px-4 py-4 whitespace-nowrap">
              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded ${
                  stage.niveau_scolaire === '2de' ? 'bg-green-100 text-green-800' :
                      stage.niveau_scolaire === '1re' ? 'bg-yellow-100 text-yellow-800' :
                          stage.niveau_scolaire === 'Tle' ? 'bg-red-100 text-red-800' :
                              'bg-gray-100 text-gray-800'
              }`}>
                {stage.niveau_scolaire || 'Non défini'}
              </span>
                      </td>
                      <td className="px-4 py-4">
                        <div className="text-sm text-gray-900 max-w-xs truncate" title={stage.famille_metiers}>
                          {stage.famille_metiers || 'Non renseigné'}
                        </div>
                      </td>
                      <td className="px-4 py-4 whitespace-nowrap">
                        <div className="text-sm text-gray-900">
                <span className={`font-medium ${
                    stage.capacity_filled >= stage.capacity_total ? 'text-red-600' :
                        stage.placesDisponibles <= 1 ? 'text-orange-600' :
                            'text-green-600'
                }`}>
                  {stage.capacity_filled}/{stage.capacity_total}
                </span>
                          <span className="text-gray-500 ml-1 text-xs block">
                  ({stage.placesDisponibles || 0} libres)
                </span>
                        </div>
                        {/* Barre de progression */}
                        <div className="mt-1 w-full bg-gray-200 rounded-full h-1.5">
                          <div
                              className={`h-1.5 rounded-full ${
                                  stage.capacity_filled >= stage.capacity_total ? 'bg-red-500' :
                                      stage.placesDisponibles <= 1 ? 'bg-orange-500' :
                                          'bg-green-500'
                              }`}
                              style={{
                                width: `${stage.capacity_total > 0 ? (stage.capacity_filled / stage.capacity_total) * 100 : 0}%`
                              }}
                          ></div>
                        </div>
                      </td>
                      <td className="px-4 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <div className="flex justify-end space-x-2">
                          <button
                              onClick={() => editStage(stage)}
                              className="text-blue-600 hover:text-blue-900 text-sm"
                              title="Modifier le stage"
                          >
                            Modifier
                          </button>
                          <button
                              onClick={() => handleDeleteStage(stage.id)}
                              className="text-red-600 hover:text-red-900 text-sm"
                              title="Supprimer le stage"
                          >
                            Supprimer
                          </button>
                        </div>
                      </td>
                    </tr>
                ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>

        {/* Modals */}
        <StageModal
            stage={editingStage}
            isOpen={isModalOpen}
            onClose={() => {
              setIsModalOpen(false)
              setEditingStage(null)
            }}
            onSuccess={handleStageSuccess}
            isNew={isNewStage}
        />
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
