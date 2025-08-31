'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { isAuthenticated, validateCredentials, storeAuth, clearAuth } from '@/lib/auth'
import LoginForm from '@/components/LoginForm'
import StageModal from '@/components/StageModal'
import type { Stage } from '@/types/stage'

export default function AdminPage() {
  const router = useRouter()
  const [stages, setStages] = useState<Stage[]>([])
  const [editingStage, setEditingStage] = useState<Stage | null>(null)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [isNewStage, setIsNewStage] = useState(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [authenticated, setAuthenticated] = useState(false)

  useEffect(() => {
    const checkAuth = () => {
      if (isAuthenticated()) {
        setAuthenticated(true)
        loadStages()
      } else {
        setAuthenticated(false)
        setLoading(false)
      }
    }
    
    checkAuth()
  }, [])

  const handleLogin = (username: string, password: string): boolean => {
    if (validateCredentials(username, password)) {
      storeAuth({
        username: username,
        isAuthenticated: true
      })
      setAuthenticated(true)
      loadStages()
      return true
    }
    return false
  }

  const handleLogout = () => {
    clearAuth()
    setAuthenticated(false)
    setStages([])
    router.push('/')
  }

  const loadStages = async () => {
    try {
      setLoading(true)
      const response = await fetch('/api/stages')
      if (!response.ok) throw new Error('Erreur de chargement')

      const data = await response.json()
      setStages(data.stages || [])
    } catch (err) {
      setError('Impossible de charger les stages')
      console.error('Erreur:', err)
    } finally {
      setLoading(false)
    }
  }

  if (!authenticated) {
    return <LoginForm onLogin={handleLogin} />
  }

  const saveStages = async (newStages: Stage[]) => {
    try {
      const response = await fetch('/api/stages', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ stages: newStages })
      })

      if (!response.ok) throw new Error('Erreur de sauvegarde')

      setStages(newStages)
      setError('')
      alert('Stage sauvegardé avec succès!')
    } catch (err) {
      setError('Impossible de sauvegarder')
      console.error('Erreur:', err)
    }
  }

  const handleAdd = () => {
    const newStage: Stage = {
      id: Date.now(),
      poste: '',
      adresse: '',
      lat: -21.1151,
      lng: 55.5364,
      placesDisponibles: 0,
      entreprise: '',
      filiere: 'CCST',
      sector: '',
      commune: '',
      capacity_total: 0,
      capacity_filled: 0,
      period: ''
    }
    setEditingStage(newStage)
    setIsNewStage(true)
    setIsModalOpen(true)
  }

  const handleEdit = (stage: Stage) => {
    setEditingStage({ ...stage })
    setIsNewStage(false)
    setIsModalOpen(true)
  }

  const handleDelete = (id: number) => {
    if (confirm('Êtes-vous sûr de vouloir supprimer ce stage ?')) {
      const newStages = stages.filter(s => s.id !== id)
      saveStages(newStages)
    }
  }

  const handleModalSave = async (stage: Stage) => {
    const existingIndex = stages.findIndex(s => s.id === stage.id)
    let newStages: Stage[]

    if (existingIndex >= 0) {
      // Modification
      newStages = [...stages]
      newStages[existingIndex] = stage
    } else {
      // Ajout
      newStages = [...stages, stage]
    }

    await saveStages(newStages)
    handleModalClose()
  }

  const handleModalClose = () => {
    setIsModalOpen(false)
    setEditingStage(null)
    setIsNewStage(false)
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg">Chargement...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="bg-white border-b border-gray-200 px-6 py-4">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Administration des Stages</h1>
            <p className="text-gray-600">{stages.length} stage(s) enregistré(s)</p>
          </div>
          <div className="flex gap-3">
            <a
            href="/"
            className="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-purple-600 transition-colors">
            Retour à la carte
            </a>
            <button
              onClick={handleAdd}
              className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
            >
              ➕ Ajouter un stage
            </button>
            <button
              onClick={handleLogout}
              className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
            >
              🚪 Déconnexion
            </button>
          </div>
        </div>
      </div>

      {error && (
        <div className="mx-6 mt-4 p-4 bg-red-50 border border-red-200 rounded-md">
          <p className="text-red-800">{error}</p>
        </div>
      )}

      <div className="p-6">
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Entreprise / Poste
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Adresse
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Filière
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Capacité
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {stages.map((stage) => (
                  <tr key={stage.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div>
                        <div className="text-sm font-medium text-gray-900">{stage.entreprise}</div>
                        <div className="text-sm text-gray-500">{stage.poste}</div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="text-sm text-gray-900">{stage.adresse}</div>
                      <div className="text-sm text-gray-500">{stage.commune}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800">
                        {stage.filiere}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      <div>Total: {stage.capacity_total}</div>
                      <div>Occupées: {stage.capacity_filled}</div>
                      <div className="font-medium text-green-600">
                        Libres: {stage.placesDisponibles}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <button
                        onClick={() => handleEdit(stage)}
                        className="text-blue-600 hover:text-blue-900 mr-3"
                      >
                        ✏️ Modifier
                      </button>
                      <button
                        onClick={() => handleDelete(stage.id)}
                        className="text-red-600 hover:text-red-900"
                      >
                        🗑️ Supprimer
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>

            {stages.length === 0 && (
              <div className="text-center py-12">
                <p className="text-gray-500">Aucun stage enregistré</p>
                <button
                  onClick={handleAdd}
                  className="mt-2 text-blue-600 hover:text-blue-800"
                >
                  Ajouter le premier stage →
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Modal séparée */}
      <StageModal
        stage={editingStage}
        isOpen={isModalOpen}
        onClose={handleModalClose}
        onSave={handleModalSave}
        isNew={isNewStage}
      />
    </div>
  )
}
