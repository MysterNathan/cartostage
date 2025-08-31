'use client'

import { useStages } from '@/contexts/StageContext'

export default function StageSelector() {
  const { filteredStages, loading, error, filters, setFilters } = useStages()

  if (loading) return <div className="p-4">Chargement des stages...</div>
  if (error) return <div className="p-4 text-red-600">Erreur: {error}</div>

  // Extraire les options uniques pour les filtres
  const communesUniques = [...new Set(filteredStages.map(s => s.commune))].sort()
  const filieresUniques = [...new Set(filteredStages.map(s => s.filiere))].sort()
  const sectorsUniques = [...new Set(filteredStages.map(s => s.sector))].sort()

  return (
      <div className="space-y-4 p-4">
        {/* Filtres */}
        <div className="flex flex-wrap gap-4">
          <select
              value={filters.commune || ''}
              onChange={(e) => setFilters({ commune: e.target.value || undefined })}
              className="border rounded px-3 py-2"
          >
            <option value="">Toutes les communes</option>
            {communesUniques.map(commune => (
                <option key={commune} value={commune}>{commune}</option>
            ))}
          </select>

          <select
              value={filters.filiere || ''}
              onChange={(e) => setFilters({ filiere: e.target.value || undefined })}
              className="border rounded px-3 py-2"
          >
            <option value="">Toutes les filières</option>
            {filieresUniques.map(filiere => (
                <option key={filiere} value={filiere}>{filiere}</option>
            ))}
          </select>

          <select
              value={filters.sector || ''}
              onChange={(e) => setFilters({ sector: e.target.value || undefined })}
              className="border rounded px-3 py-2"
          >
            <option value="">Tous les secteurs</option>
            {sectorsUniques.map(sector => (
                <option key={sector} value={sector}>{sector}</option>
            ))}
          </select>
        </div>

        {/* Liste des stages */}
        <div className="space-y-2">
          <h3 className="font-semibold">
            {filteredStages.length} stage(s) trouvé(s)
          </h3>

          {filteredStages.map(stage => (
              <div key={stage.id} className="border rounded p-3">
                <h4 className="font-medium">{stage.poste}</h4>
                <p className="text-sm text-gray-600">{stage.entreprise} - {stage.commune}</p>
                <p className="text-sm">
                  Places: {stage.capacity_filled}/{stage.capacity_total}
                </p>
              </div>
          ))}
        </div>
      </div>
  )
}
