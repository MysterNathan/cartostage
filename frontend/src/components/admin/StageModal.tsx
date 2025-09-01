'use client'

import { useState, useEffect } from 'react'
import type { Stage } from '@/types/stage'
import { addStage, updateStage } from '@/lib/stageApi'
import FiliereSelect from '@/components/admin/FiliereSelect'

interface StageModalProps {
  stage: Stage | null
  isOpen: boolean
  onClose: () => void
  onSave: (stage: Stage) => void
  isNew?: boolean
}

export default function StageModal({ stage, isOpen, onClose, onSave, isNew = false }: StageModalProps) {
  const [editingStage, setEditingStage] = useState<Stage | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    if (stage) {
      setEditingStage({ ...stage })
    }
  }, [stage])

  // Fermer la modal en cliquant sur l'overlay
  const handleOverlayClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose()
    }
  }

  if (!isOpen || !editingStage) return null

  const handleSave = async () => {
    if (!editingStage) return

    try {
      setIsLoading(true)
      setError('')

      // Calculer les places disponibles
      editingStage.placesDisponibles = Math.max(0, editingStage.capacity_total - editingStage.capacity_filled)

      let savedStage: Stage

      if (isNew) {
        // Création d'un nouveau stage
        const stageData = {
          entreprise: editingStage.entreprise,
          poste: editingStage.poste,
          adresse: editingStage.adresse,
          commune: editingStage.commune,
          sector: editingStage.sector,
          filiere: editingStage.filiere,
          lat: editingStage.lat,
          lng: editingStage.lng,
          capacity_total: editingStage.capacity_total,
          capacity_filled: editingStage.capacity_filled,
          placesDisponibles: editingStage.placesDisponibles,
          period: editingStage.period
        }
        savedStage = await addStage(stageData)
      } else {
        // Modification d'un stage existant
        const stageData = {
          entreprise: editingStage.entreprise,
          poste: editingStage.poste,
          adresse: editingStage.adresse,
          commune: editingStage.commune,
          sector: editingStage.sector,
          filiere: editingStage.filiere,
          lat: editingStage.lat,
          lng: editingStage.lng,
          capacity_total: editingStage.capacity_total,
          capacity_filled: editingStage.capacity_filled,
          placesDisponibles: editingStage.placesDisponibles,
          period: editingStage.period
        }
        savedStage = await updateStage(editingStage.id, stageData)
      }

      // Notifier le composant parent du succès
      onSave(savedStage)

    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Erreur lors de la sauvegarde'
      setError(errorMessage)
      console.error('Erreur:', err)
    } finally {
      setIsLoading(false)
    }
  }

  return (
      <div
          className="fixed inset-0 flex items-center justify-center p-4 z-50"
          onClick={handleOverlayClick}
      >
        {/* Overlay transparent avec légère teinte */}
        <div className="absolute inset-0 backdrop-blur-sm"></div>

        {/* Modal */}
        <div className="relative bg-white rounded-lg w-full max-w-2xl max-h-[90vh] overflow-hidden shadow-2xl border border-gray-200">
          {/* Header avec bouton fermer */}
          <div className="sticky top-0 bg-white p-6 border-b border-gray-200 flex justify-between items-center">
            <h2 className="text-xl font-semibold text-gray-900">
              {isNew ? 'Ajouter un stage' : 'Modifier le stage'}
            </h2>
            <button
                onClick={onClose}
                disabled={isLoading}
                className="p-1 hover:bg-gray-100 rounded-full transition-colors disabled:opacity-50"
                aria-label="Fermer"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          {/* Message d'erreur */}
          {error && (
              <div className="mx-6 mt-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-md">
                {error}
              </div>
          )}

          {/* Contenu scrollable */}
          <div className="overflow-y-auto max-h-[calc(90vh-140px)]">
            <div className="p-6 space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Entreprise *
                  </label>
                  <input
                      type="text"
                      value={editingStage.entreprise}
                      onChange={(e) => setEditingStage({ ...editingStage, entreprise: e.target.value })}
                      disabled={isLoading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
                      required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Poste *
                  </label>
                  <input
                      type="text"
                      value={editingStage.poste}
                      onChange={(e) => setEditingStage({ ...editingStage, poste: e.target.value })}
                      disabled={isLoading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
                      required
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Adresse complète *
                </label>
                <input
                    type="text"
                    value={editingStage.adresse}
                    onChange={(e) => setEditingStage({ ...editingStage, adresse: e.target.value })}
                    placeholder="123 Rue de la République, Saint-Denis, La Réunion"
                    disabled={isLoading}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
                    required
                />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Commune
                  </label>
                  <input
                      type="text"
                      value={editingStage.commune || ''}
                      onChange={(e) => setEditingStage({ ...editingStage, commune: e.target.value })}
                      disabled={isLoading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Secteur
                  </label>
                  <input
                      type="text"
                      value={editingStage.sector || ''}
                      onChange={(e) => setEditingStage({ ...editingStage, sector: e.target.value })}
                      disabled={isLoading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
                  />
                </div>
              </div>

              {/* Composant FiliereSelect - utilise maintenant l'API des filières */}
              <div>
                <FiliereSelect
                    value={editingStage.filiere}
                    onChange={(code) => setEditingStage({ ...editingStage, filiere: code })}
                    required
                    allowCreate={false}
                    disabled={isLoading}
                />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Latitude
                  </label>
                  <input
                      type="number"
                      step="0.000001"
                      value={editingStage.lat}
                      onChange={(e) => setEditingStage({ ...editingStage, lat: parseFloat(e.target.value) || 0 })}
                      disabled={isLoading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Longitude
                  </label>
                  <input
                      type="number"
                      step="0.000001"
                      value={editingStage.lng}
                      onChange={(e) => setEditingStage({ ...editingStage, lng: parseFloat(e.target.value) || 0 })}
                      disabled={isLoading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Capacité totale *
                  </label>
                  <input
                      type="number"
                      min="0"
                      value={editingStage.capacity_total}
                      onChange={(e) => {
                        const total = parseInt(e.target.value) || 0
                        setEditingStage({
                          ...editingStage,
                          capacity_total: total,
                          capacity_filled: Math.min(editingStage.capacity_filled, total),
                          placesDisponibles: total - Math.min(editingStage.capacity_filled, total)
                        })
                      }}
                      disabled={isLoading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
                      required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Places occupées
                  </label>
                  <input
                      type="number"
                      min="0"
                      max={editingStage.capacity_total}
                      value={editingStage.capacity_filled}
                      onChange={(e) => {
                        const filled = Math.min(parseInt(e.target.value) || 0, editingStage.capacity_total)
                        setEditingStage({
                          ...editingStage,
                          capacity_filled: filled,
                          placesDisponibles: editingStage.capacity_total - filled
                        })
                      }}
                      disabled={isLoading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Périodes
                </label>
                <input
                    type="text"
                    value={editingStage.period || ''}
                    onChange={(e) => setEditingStage({ ...editingStage, period: e.target.value })}
                    placeholder="Ex: Oct–Nov; Mar–Avr"
                    disabled={isLoading}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100"
                />
              </div>

              <div className="bg-gray-50 p-3 rounded-md">
                <div className="text-sm text-gray-600">
                  <strong>Places libres calculées:</strong> {editingStage.placesDisponibles}
                </div>
              </div>
            </div>
          </div>

          {/* Footer avec boutons */}
          <div className="sticky bottom-0 bg-white p-6 border-t border-gray-200 flex justify-end gap-3">
            <button
                onClick={onClose}
                disabled={isLoading}
                className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50"
            >
              Annuler
            </button>
            <button
                onClick={handleSave}
                disabled={isLoading}
                className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50 flex items-center"
            >
              {isLoading && (
                  <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
              )}
              {isLoading ? 'Sauvegarde...' : 'Sauvegarder'}
            </button>
          </div>
        </div>
      </div>
  )
}
