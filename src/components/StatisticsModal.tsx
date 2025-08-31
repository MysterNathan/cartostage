'use client'

import { useMemo } from 'react'
import type { Stage } from '@/types/stage'

interface StatisticsModalProps {
  stages: Stage[]
  isOpen: boolean
  onClose: () => void
}

export default function StatisticsModal({ stages, isOpen, onClose }: StatisticsModalProps) {
  const stats = useMemo(() => {
    if (!stages || stages.length === 0) {
      return {
        totalStages: 0,
        totalCapacity: 0,
        totalFilled: 0,
        totalAvailable: 0,
        fillRate: 0,
        byFiliere: {},
        byCommune: {},
        bySector: {}
      }
    }

    const totalStages = stages.length
    const totalCapacity = stages.reduce((sum, stage) => sum + stage.capacity_total, 0)
    const totalFilled = stages.reduce((sum, stage) => sum + stage.capacity_filled, 0)
    const totalAvailable = totalCapacity - totalFilled
    const fillRate = totalCapacity > 0 ? (totalFilled / totalCapacity) * 100 : 0

    // Statistiques par filière
    const byFiliere = stages.reduce((acc, stage) => {
      if (!acc[stage.filiere]) {
        acc[stage.filiere] = {
          count: 0,
          capacity: 0,
          filled: 0,
          available: 0
        }
      }
      acc[stage.filiere].count += 1
      acc[stage.filiere].capacity += stage.capacity_total
      acc[stage.filiere].filled += stage.capacity_filled
      acc[stage.filiere].available += stage.placesDisponibles || 0
      return acc
    }, {} as Record<string, any>)

    // Statistiques par commune
    const byCommune = stages.reduce((acc, stage) => {
      const commune = stage.commune || 'Non renseigné'
      if (!acc[commune]) {
        acc[commune] = {
          count: 0,
          capacity: 0,
          filled: 0,
          available: 0
        }
      }
      acc[commune].count += 1
      acc[commune].capacity += stage.capacity_total
      acc[commune].filled += stage.capacity_filled
      acc[commune].available += stage.placesDisponibles || 0
      return acc
    }, {} as Record<string, any>)

    // Statistiques par secteur
    const bySector = stages.reduce((acc, stage) => {
      const sector = stage.sector || 'Non renseigné'
      if (!acc[sector]) {
        acc[sector] = {
          count: 0,
          capacity: 0,
          filled: 0,
          available: 0
        }
      }
      acc[sector].count += 1
      acc[sector].capacity += stage.capacity_total
      acc[sector].filled += stage.capacity_filled
      acc[sector].available += stage.placesDisponibles || 0
      return acc
    }, {} as Record<string, any>)

    return {
      totalStages,
      totalCapacity,
      totalFilled,
      totalAvailable,
      fillRate,
      byFiliere,
      byCommune,
      bySector
    }
  }, [stages])

  const handleOverlayClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose()
    }
  }

  const handleModalClick = (e: React.MouseEvent) => {
    e.stopPropagation()
  }

  if (!isOpen) return null

  return (
    <div 
      className="fixed inset-0 flex items-center justify-center p-4 z-50 cursor-pointer"
      onClick={handleOverlayClick}
    >
      <div className="absolute inset-0 backdrop-blur-sm bg-black bg-opacity-10"></div>

      <div 
        className="relative bg-white rounded-lg w-full max-w-4xl max-h-[90vh] overflow-hidden shadow-2xl border border-gray-200 cursor-default"
        onClick={handleModalClick}
      >
        {/* Header */}
        <div className="sticky top-0 bg-white p-6 border-b border-gray-200 flex justify-between items-center">
          <h2 className="text-xl font-semibold text-gray-900 flex items-center">
            <svg className="w-6 h-6 mr-2 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            Statistiques des stages
          </h2>
          <button
            onClick={onClose}
            className="p-1 hover:bg-gray-100 rounded-full transition-colors"
          >
            <svg className="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Contenu scrollable */}
        <div className="overflow-y-auto max-h-[calc(90vh-120px)] p-6">
          {/* Vue d'ensemble */}
          <div className="mb-8">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">Vue d'ensemble</h3>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                <div className="text-2xl font-bold text-blue-600">{stats.totalStages}</div>
                <div className="text-sm text-blue-800">Stages totaux</div>
              </div>
              <div className="bg-green-50 p-4 rounded-lg border border-green-200">
                <div className="text-2xl font-bold text-green-600">{stats.totalCapacity}</div>
                <div className="text-sm text-green-800">Places totales</div>
              </div>
              <div className="bg-orange-50 p-4 rounded-lg border border-orange-200">
                <div className="text-2xl font-bold text-orange-600">{stats.totalFilled}</div>
                <div className="text-sm text-orange-800">Places occupées</div>
              </div>
              <div className="bg-purple-50 p-4 rounded-lg border border-purple-200">
                <div className="text-2xl font-bold text-purple-600">{stats.totalAvailable}</div>
                <div className="text-sm text-purple-800">Places libres</div>
              </div>
            </div>
            
            {/* Taux de remplissage */}
            <div className="mt-4 bg-gray-50 p-4 rounded-lg">
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm font-medium text-gray-700">Taux de remplissage</span>
                <span className="text-sm font-semibold text-gray-900">{stats.fillRate.toFixed(1)}%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-3">
                <div 
                  className="bg-gradient-to-r from-blue-500 to-blue-600 h-3 rounded-full transition-all duration-300"
                  style={{ width: `${Math.min(stats.fillRate, 100)}%` }}
                ></div>
              </div>
            </div>
          </div>

          {/* Par filière */}
          <div className="mb-8">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">Par filière</h3>
            <div className="space-y-3">
              {Object.entries(stats.byFiliere).map(([filiere, data]: [string, any]) => (
                <div key={filiere} className="bg-gray-50 p-4 rounded-lg">
                  <div className="flex justify-between items-center mb-2">
                    <h4 className="font-medium text-gray-800">{filiere}</h4>
                    <span className="text-sm text-gray-600">{data.count} stages</span>
                  </div>
                  <div className="grid grid-cols-3 gap-4 text-sm">
                    <div>
                      <span className="text-gray-600">Capacité: </span>
                      <span className="font-medium">{data.capacity}</span>
                    </div>
                    <div>
                      <span className="text-gray-600">Occupées: </span>
                      <span className="font-medium text-orange-600">{data.filled}</span>
                    </div>
                    <div>
                      <span className="text-gray-600">Libres: </span>
                      <span className="font-medium text-green-600">{data.available}</span>
                    </div>
                  </div>
                  <div className="mt-2">
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div 
                        className="bg-gradient-to-r from-orange-400 to-orange-500 h-2 rounded-full"
                        style={{ width: `${data.capacity > 0 ? (data.filled / data.capacity) * 100 : 0}%` }}
                      ></div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Par commune */}
          <div className="mb-8">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">Par commune</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {Object.entries(stats.byCommune)
                .sort(([,a]: [string, any], [,b]: [string, any]) => b.count - a.count)
                .slice(0, 10)
                .map(([commune, data]: [string, any]) => (
                <div key={commune} className="bg-gray-50 p-3 rounded-lg">
                  <div className="flex justify-between items-center mb-1">
                    <h4 className="font-medium text-gray-800 text-sm">{commune}</h4>
                    <span className="text-xs text-gray-600">{data.count} stages</span>
                  </div>
                  <div className="flex justify-between text-xs text-gray-600">
                    <span>Cap: {data.capacity}</span>
                    <span>Occ: {data.filled}</span>
                    <span>Lib: {data.available}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Par secteur */}
          <div>
            <h3 className="text-lg font-semibold text-gray-800 mb-4">Par secteur d'activité</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {Object.entries(stats.bySector)
                .sort(([,a]: [string, any], [,b]: [string, any]) => b.count - a.count)
                .slice(0, 10)
                .map(([sector, data]: [string, any]) => (
                <div key={sector} className="bg-gray-50 p-3 rounded-lg">
                  <div className="flex justify-between items-center mb-1">
                    <h4 className="font-medium text-gray-800 text-sm">{sector}</h4>
                    <span className="text-xs text-gray-600">{data.count} stages</span>
                  </div>
                  <div className="flex justify-between text-xs text-gray-600">
                    <span>Cap: {data.capacity}</span>
                    <span>Occ: {data.filled}</span>
                    <span>Lib: {data.available}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="sticky bottom-0 bg-white p-4 border-t border-gray-200 flex justify-end">
          <button
            onClick={onClose}
            className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
          >
            Fermer
          </button>
        </div>
      </div>
    </div>
  )
}
