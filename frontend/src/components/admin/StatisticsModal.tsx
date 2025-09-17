'use client'

import { useMemo, useState } from 'react'
import type { Stage } from '@/types/stage'

interface StatisticsModalProps {
  stages: Stage[]
  isOpen: boolean
  onClose: () => void
}

type TabType = 'general' | 'parcours-y'

export default function StatisticsModal({ stages, isOpen, onClose }: StatisticsModalProps) {
  const [activeTab, setActiveTab] = useState<TabType>('general')

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
        bySector: {},
        byParcours: {},
        byFamilleMetiers: {},
        byNiveauScolaire: {},
        // Stats spécifiques au Parcours Y
        parcoursY: {
          totalStages: 0,
          totalCapacity: 0,
          totalFilled: 0,
          totalAvailable: 0,
          fillRate: 0,
          byFiliere: {},
          byPeriod: {},
          byCommune: {},
          bySector: {},
          byEntreprise: {}
        }
      }
    }

    const totalStages = stages.length
    const totalCapacity = stages.reduce((sum, stage) => sum + stage.capacity_total, 0)
    const totalFilled = stages.reduce((sum, stage) => sum + stage.capacity_filled, 0)
    const totalAvailable = totalCapacity - totalFilled
    const fillRate = totalCapacity > 0 ? (totalFilled / totalCapacity) * 100 : 0

    // Fonction utilitaire pour créer les stats
    const createStats = (key: keyof Stage, stagesList = stages) => {
      return stagesList.reduce((acc, stage) => {
        const value = stage[key] || 'Non renseigné'
        if (!acc[value]) {
          acc[value] = {
            count: 0,
            capacity: 0,
            filled: 0,
            available: 0
          }
        }
        acc[value].count += 1
        acc[value].capacity += stage.capacity_total
        acc[value].filled += stage.capacity_filled
        acc[value].available += stage.places_disponibles || 0
        return acc
      }, {} as Record<string, any>)
    }

    // Statistiques générales
    const byFiliere = createStats('filiere')
    const byCommune = createStats('commune')
    const bySector = createStats('sector')
    const byParcours = createStats('parcours')
    const byFamilleMetiers = createStats('famille_metiers')
    const byNiveauScolaire = createStats('niveau_scolaire')

    // Filtrer les stages du Parcours Y (ici j'assume que c'est identifiable par certains critères)
    // Vous devrez adapter cette logique selon vos données
    const parcoursYStages = stages.filter(stage =>
        stage.parcours === 'scolaire' &&
        (stage.niveauScolaire === '2de' || stage.period?.includes('découverte'))
    )

    const parcoursYStats = {
      totalStages: parcoursYStages.length,
      totalCapacity: parcoursYStages.reduce((sum, stage) => sum + stage.capacity_total, 0),
      totalFilled: parcoursYStages.reduce((sum, stage) => sum + stage.capacity_filled, 0),
      totalAvailable: parcoursYStages.reduce((sum, stage) => sum + (stage.places_disponibles || 0), 0),
      fillRate: 0,
      byFiliere: createStats('filiere', parcoursYStages),
      byPeriod: createStats('period', parcoursYStages),
      byCommune: createStats('commune', parcoursYStages),
      bySector: createStats('sector', parcoursYStages),
      byEntreprise: createStats('enterprise', parcoursYStages)
    }

    parcoursYStats.fillRate = parcoursYStats.totalCapacity > 0
        ? (parcoursYStats.totalFilled / parcoursYStats.totalCapacity) * 100
        : 0

    return {
      totalStages,
      totalCapacity,
      totalFilled,
      totalAvailable,
      fillRate,
      byFiliere,
      byCommune,
      bySector,
      byParcours,
      byFamilleMetiers,
      byNiveauScolaire,
      parcoursY: parcoursYStats
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

  // Composant réutilisable pour les sections de statistiques
  const StatsSection = ({
                          title,
                          data,
                          columns = 1,
                          showProgress = true,
                          limit = 10,
                          bgColor = "bg-white",
                          borderColor = "border-gray-300"
                        }: {
    title: string
    data: Record<string, any>
    columns?: 1 | 2
    showProgress?: boolean
    limit?: number
    bgColor?: string
    borderColor?: string
  }) => (
      <div className="mb-8">
        <h3 className="text-lg font-semibold text-gray-800 mb-4">{title}</h3>
        <div className={`grid grid-cols-1 ${columns === 2 ? 'md:grid-cols-2' : ''} gap-4`}>
          {Object.entries(data)
              .sort(([,a]: [string, any], [,b]: [string, any]) => b.count - a.count)
              .slice(0, limit)
              .map(([key, data]: [string, any]) => (
                  <div key={key} className={`${bgColor} p-4 rounded-lg border ${borderColor}`}>
                    <div className="flex justify-between items-center mb-2">
                      <h4 className="font-medium text-gray-900">{key}</h4>
                      <span className="text-sm text-gray-700 bg-gray-100 px-2 py-1 rounded">{data.count} stages</span>
                    </div>
                    <div className="grid grid-cols-3 gap-4 text-sm">
                      <div>
                        <span className="text-gray-600">Capacité: </span>
                        <span className="font-semibold text-gray-900">{data.capacity}</span>
                      </div>
                      <div>
                        <span className="text-gray-600">Occupées: </span>
                        <span className="font-semibold text-orange-700">{data.filled}</span>
                      </div>
                      <div>
                        <span className="text-gray-600">Libres: </span>
                        <span className="font-semibold text-green-700">{data.available}</span>
                      </div>
                    </div>
                    {showProgress && (
                        <div className="mt-3">
                          <div className="w-full bg-gray-200 rounded-full h-2">
                            <div
                                className="bg-gradient-to-r from-orange-500 to-orange-600 h-2 rounded-full"
                                style={{ width: `${data.capacity > 0 ? (data.filled / data.capacity) * 100 : 0}%` }}
                            ></div>
                          </div>
                        </div>
                    )}
                  </div>
              ))}
        </div>
      </div>
  )

  if (!isOpen) return null

  return (
      <div
          className="fixed inset-0 flex items-center justify-center p-4 z-50 cursor-pointer"
          onClick={handleOverlayClick}
      >
        <div className="absolute inset-0 backdrop-blur-sm bg-black bg-opacity-30"></div>

        <div
            className="relative bg-white rounded-lg w-full max-w-6xl max-h-[90vh] overflow-hidden shadow-2xl border border-gray-300 cursor-default"
            onClick={handleModalClick}
        >
          {/* Header */}
          <div className="sticky top-0 bg-white p-6 border-b border-gray-300 flex justify-between items-center">
            <h2 className="text-xl font-semibold text-gray-900 flex items-center">
              <svg className="w-6 h-6 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
              Statistiques des stages
            </h2>
            <button
                onClick={onClose}
                className="p-1 hover:bg-gray-100 rounded-full transition-colors"
            >
              <svg className="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          {/* Tabs */}
          <div className="bg-gray-50 border-b border-gray-300">
            <div className="flex">
              <button
                  onClick={() => setActiveTab('general')}
                  className={`px-6 py-3 font-medium text-sm border-b-2 transition-colors ${
                      activeTab === 'general'
                          ? 'border-blue-500 text-blue-600 bg-white'
                          : 'border-transparent text-gray-600 hover:text-gray-800 hover:bg-gray-100'
                  }`}
              >
                Statistiques générales
              </button>
              <button
                  onClick={() => setActiveTab('parcours-y')}
                  className={`px-6 py-3 font-medium text-sm border-b-2 transition-colors ${
                      activeTab === 'parcours-y'
                          ? 'border-blue-500 text-blue-600 bg-white'
                          : 'border-transparent text-gray-600 hover:text-gray-800 hover:bg-gray-100'
                  }`}
              >
                Parcours Y
              </button>
            </div>
          </div>

          {/* Contenu scrollable */}
          <div className="overflow-y-auto max-h-[calc(90vh-180px)] p-6">
            {activeTab === 'general' && (
                <>
                  {/* Vue d'ensemble générale */}
                  <div className="mb-8">
                    <h3 className="text-lg font-semibold text-gray-800 mb-4">Vue d'ensemble</h3>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                      <div className="bg-blue-100 p-4 rounded-lg border border-blue-300">
                        <div className="text-2xl font-bold text-blue-800">{stats.totalStages}</div>
                        <div className="text-sm font-medium text-blue-700">Stages totaux</div>
                      </div>
                      <div className="bg-green-100 p-4 rounded-lg border border-green-300">
                        <div className="text-2xl font-bold text-green-800">{stats.totalCapacity}</div>
                        <div className="text-sm font-medium text-green-700">Places totales</div>
                      </div>
                      <div className="bg-orange-100 p-4 rounded-lg border border-orange-300">
                        <div className="text-2xl font-bold text-orange-800">{stats.totalFilled}</div>
                        <div className="text-sm font-medium text-orange-700">Places occupées</div>
                      </div>
                      <div className="bg-purple-100 p-4 rounded-lg border border-purple-300">
                        <div className="text-2xl font-bold text-purple-800">{stats.totalAvailable}</div>
                        <div className="text-sm font-medium text-purple-700">Places libres</div>
                      </div>
                    </div>

                    {/* Taux de remplissage */}
                    <div className="mt-4 bg-gray-100 p-4 rounded-lg border border-gray-300">
                      <div className="flex justify-between items-center mb-2">
                        <span className="text-sm font-semibold text-gray-800">Taux de remplissage</span>
                        <span className="text-sm font-bold text-gray-900">{stats.fillRate.toFixed(1)}%</span>
                      </div>
                      <div className="w-full bg-gray-300 rounded-full h-3">
                        <div
                            className="bg-gradient-to-r from-blue-600 to-blue-700 h-3 rounded-full transition-all duration-300"
                            style={{ width: `${Math.min(stats.fillRate, 100)}%` }}
                        ></div>
                      </div>
                    </div>
                  </div>

                  {/* Statistiques par parcours */}
                  <div className="mb-8">
                    <h3 className="text-lg font-semibold text-gray-800 mb-4">Par type de parcours</h3>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      {Object.entries(stats.byParcours).map(([parcours, data]: [string, any]) => {
                        const colorConfigs = {
                          scolaire: { bg: 'bg-blue-100', border: 'border-blue-300', text: 'text-blue-800' },
                          apprentissage: { bg: 'bg-green-100', border: 'border-green-300', text: 'text-green-800' },
                          mixte: { bg: 'bg-purple-100', border: 'border-purple-300', text: 'text-purple-800' }
                        }
                        const config = colorConfigs[parcours as keyof typeof colorConfigs] ||
                            { bg: 'bg-gray-100', border: 'border-gray-300', text: 'text-gray-800' }

                        return (
                            <div key={parcours} className={`${config.bg} p-4 rounded-lg border ${config.border}`}>
                              <div className="text-center mb-3">
                                <div className={`text-2xl font-bold ${config.text} capitalize`}>{parcours}</div>
                                <div className={`text-sm font-medium ${config.text}`}>{data.count} stages</div>
                              </div>
                              <div className="space-y-2 text-sm">
                                <div className="flex justify-between">
                                  <span className="text-gray-700">Capacité:</span>
                                  <span className="font-semibold text-gray-900">{data.capacity}</span>
                                </div>
                                <div className="flex justify-between">
                                  <span className="text-gray-700">Occupées:</span>
                                  <span className="font-semibold text-gray-900">{data.filled}</span>
                                </div>
                                <div className="flex justify-between">
                                  <span className="text-gray-700">Libres:</span>
                                  <span className="font-semibold text-gray-900">{data.available}</span>
                                </div>
                              </div>
                            </div>
                        )
                      })}
                    </div>
                  </div>

                  {/* Statistiques par niveau scolaire */}
                  <div className="mb-8">
                    <h3 className="text-lg font-semibold text-gray-800 mb-4">Par niveau scolaire</h3>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      {(['2de', '1re', 'Tle'] as const).map((niveau) => {
                        const data = stats.byNiveauScolaire[niveau]
                        if (!data) return null

                        const colorConfigs = {
                          '2de': { bg: 'bg-emerald-100', border: 'border-emerald-300', text: 'text-emerald-800' },
                          '1re': { bg: 'bg-amber-100', border: 'border-amber-300', text: 'text-amber-800' },
                          'Tle': { bg: 'bg-red-100', border: 'border-red-300', text: 'text-red-800' }
                        }
                        const config = colorConfigs[niveau]

                        return (
                            <div key={niveau} className={`${config.bg} p-4 rounded-lg border ${config.border}`}>
                              <div className="text-center mb-3">
                                <div className={`text-2xl font-bold ${config.text}`}>{niveau}</div>
                                <div className={`text-sm font-medium ${config.text}`}>{data.count} stages</div>
                              </div>
                              <div className="space-y-2 text-sm">
                                <div className="flex justify-between">
                                  <span className="text-gray-700">Capacité:</span>
                                  <span className="font-semibold text-gray-900">{data.capacity}</span>
                                </div>
                                <div className="flex justify-between">
                                  <span className="text-gray-700">Occupées:</span>
                                  <span className="font-semibold text-gray-900">{data.filled}</span>
                                </div>
                                <div className="flex justify-between">
                                  <span className="text-gray-700">Libres:</span>
                                  <span className="font-semibold text-gray-900">{data.available}</span>
                                </div>
                              </div>
                            </div>
                        )
                      })}
                    </div>
                  </div>

                  {/* Par famille de métiers */}
                  <StatsSection
                      title="Par famille de métiers"
                      data={stats.byFamilleMetiers}
                      columns={2}
                      bgColor="bg-indigo-50"
                      borderColor="border-indigo-200"
                  />

                  {/* Par filière */}
                  <StatsSection
                      title="Par filière"
                      data={stats.byFiliere}
                      columns={1}
                      bgColor="bg-orange-50"
                      borderColor="border-orange-200"
                  />

                  {/* Par commune */}
                  <StatsSection
                      title="Par commune (Top 10)"
                      data={stats.byCommune}
                      columns={2}
                      showProgress={false}
                      bgColor="bg-teal-50"
                      borderColor="border-teal-200"
                  />

                  {/* Par secteur */}
                  <StatsSection
                      title="Par secteur d'activité (Top 10)"
                      data={stats.bySector}
                      columns={2}
                      showProgress={false}
                      bgColor="bg-cyan-50"
                      borderColor="border-cyan-200"
                  />
                </>
            )}

            {activeTab === 'parcours-y' && (
                <>
                  {/* Vue d'ensemble Parcours Y */}
                  <div className="mb-8">
                    <div className="bg-gradient-to-r from-purple-100 to-pink-100 p-6 rounded-lg border border-purple-200 mb-6">
                      <h3 className="text-2xl font-bold text-purple-800 mb-2">Parcours Y - Découverte des métiers</h3>
                      <p className="text-purple-700">
                        Statistiques spécifiques aux stages de découverte professionnelle permettant aux élèves d'explorer différentes filières avant leur spécialisation.
                      </p>
                    </div>

                    <h3 className="text-lg font-semibold text-gray-800 mb-4">Vue d'ensemble Parcours Y</h3>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                      <div className="bg-purple-100 p-4 rounded-lg border border-purple-300">
                        <div className="text-2xl font-bold text-purple-800">{stats.parcoursY.totalStages}</div>
                        <div className="text-sm font-medium text-purple-700">Stages découverte</div>
                      </div>
                      <div className="bg-pink-100 p-4 rounded-lg border border-pink-300">
                        <div className="text-2xl font-bold text-pink-800">{stats.parcoursY.totalCapacity}</div>
                        <div className="text-sm font-medium text-pink-700">Places totales</div>
                      </div>
                      <div className="bg-indigo-100 p-4 rounded-lg border border-indigo-300">
                        <div className="text-2xl font-bold text-indigo-800">{stats.parcoursY.totalFilled}</div>
                        <div className="text-sm font-medium text-indigo-700">Places occupées</div>
                      </div>
                      <div className="bg-violet-100 p-4 rounded-lg border border-violet-300">
                        <div className="text-2xl font-bold text-violet-800">{stats.parcoursY.totalAvailable}</div>
                        <div className="text-sm font-medium text-violet-700">Places libres</div>
                      </div>
                    </div>

                    {/* Taux de remplissage Parcours Y */}
                    <div className="mt-4 bg-gradient-to-r from-purple-50 to-pink-50 p-4 rounded-lg border border-purple-200">
                      <div className="flex justify-between items-center mb-2">
                        <span className="text-sm font-semibold text-purple-800">Taux de remplissage Parcours Y</span>
                        <span className="text-sm font-bold text-purple-900">{stats.parcoursY.fillRate.toFixed(1)}%</span>
                      </div>
                      <div className="w-full bg-purple-200 rounded-full h-3">
                        <div
                            className="bg-gradient-to-r from-purple-600 to-pink-600 h-3 rounded-full transition-all duration-300"
                            style={{ width: `${Math.min(stats.parcoursY.fillRate, 100)}%` }}
                        ></div>
                      </div>
                    </div>
                  </div>

                  {/* Répartition par filière de découverte */}
                  <StatsSection
                      title="Filières de découverte"
                      data={stats.parcoursY.byFiliere}
                      columns={2}
                      bgColor="bg-purple-50"
                      borderColor="border-purple-200"
                  />

                  {/* Par période de stage */}
                  <StatsSection
                      title="Périodes de stage"
                      data={stats.parcoursY.byPeriod}
                      columns={2}
                      bgColor="bg-pink-50"
                      borderColor="border-pink-200"
                  />

                  {/* Top entreprises d'accueil */}
                  <StatsSection
                      title="Top entreprises d'accueil (Parcours Y)"
                      data={stats.parcoursY.byEntreprise}
                      columns={1}
                      limit={8}
                      bgColor="bg-indigo-50"
                      borderColor="border-indigo-200"
                  />

                  {/* Répartition géographique */}
                  <StatsSection
                      title="Répartition géographique (Parcours Y)"
                      data={stats.parcoursY.byCommune}
                      columns={2}
                      showProgress={false}
                      limit={8}
                      bgColor="bg-violet-50"
                      borderColor="border-violet-200"
                  />

                  {/* Diversité sectorielle */}
                  <StatsSection
                      title="Diversité des secteurs explorés"
                      data={stats.parcoursY.bySector}
                      columns={2}
                      showProgress={false}
                      bgColor="bg-cyan-50"
                      borderColor="border-cyan-200"
                  />

                  {/* Message informatif */}
                  <div className="bg-blue-50 border border-blue-200 p-4 rounded-lg">
                    <div className="flex items-start">
                      <svg className="w-5 h-5 text-blue-600 mt-0.5 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <div>
                        <h4 className="text-sm font-semibold text-blue-800 mb-1">À propos du Parcours Y</h4>
                        <p className="text-sm text-blue-700">
                          Le Parcours Y permet aux élèves de 2de professionnelle de découvrir plusieurs filières avant de choisir leur spécialisation définitive.
                          Cette approche réduit les risques d'erreur d'orientation et favorise l'épanouissement des élèves.
                        </p>
                      </div>
                    </div>
                  </div>
                </>
            )}
          </div>

          {/* Footer */}
          <div className="sticky bottom-0 bg-white p-4 border-t border-gray-300 flex justify-end">
            <button
                onClick={onClose}
                className="px-4 py-2 bg-gray-200 text-gray-800 rounded-lg hover:bg-gray-300 transition-colors font-medium"
            >
              Fermer
            </button>
          </div>
        </div>
      </div>
  )
}
