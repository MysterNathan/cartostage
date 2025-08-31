'use client'
import type { Stage } from '@/types/stage'

interface StageListProps {
  stages: Stage[]
  onStageSelect: (stage: Stage) => void
  selectedStage?: Stage | null
}

export default function StageList({ stages, onStageSelect, selectedStage }: StageListProps) {
  if (stages.length === 0) {
    return (
      <div className="bg-white p-4 rounded-lg shadow-sm border">
        <p className="text-gray-500 text-center">Aucun stage trouvé</p>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg shadow-sm border">
      <div className="p-4 border-b border-gray-200">
        <h3 className="font-semibold">Stages disponibles</h3>
        <p className="text-sm text-gray-500">{stages.length} résultat{stages.length > 1 ? 's' : ''}</p>
      </div>
      
      <div className="max-h-96 overflow-y-auto">
        {stages.map((stage) => {
          const isSelected = selectedStage?.id === stage.id
          const tauxOccupation = stage.capacity_total > 0 
            ? Math.round((stage.capacity_filled / stage.capacity_total) * 100)
            : 0

          return (
            <div
              key={stage.id}
              onClick={() => onStageSelect(stage)}
              className={`
                p-4 border-b border-gray-100 cursor-pointer transition-colors
                ${isSelected 
                  ? 'bg-blue-50 border-l-4 border-l-blue-500' 
                  : 'hover:bg-gray-50'
                }
              `}
            >
              <div className="flex items-start justify-between mb-2">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-1">
                    <span className={`
                      px-2 py-1 text-xs font-medium rounded-full
                      ${stage.filiere === 'CCST' 
                        ? 'bg-blue-100 text-blue-800' 
                        : 'bg-green-100 text-green-800'
                      }
                    `}>
                      {stage.filiere}
                    </span>
                    <span className={`
                      px-2 py-1 text-xs rounded-full font-medium
                      ${stage.placesDisponibles > 0 
                        ? 'bg-green-100 text-green-800' 
                        : 'bg-red-100 text-red-800'
                      }
                    `}>
                      {stage.placesDisponibles} place{stage.placesDisponibles !== 1 ? 's' : ''} libre{stage.placesDisponibles !== 1 ? 's' : ''}
                    </span>
                  </div>
                  
                  <h4 className="font-semibold text-gray-900 text-sm mb-1">
                    {stage.poste}
                  </h4>
                  
                  <div className="text-sm text-gray-600 mb-2">
                    <div className="font-medium">{stage.entreprise}</div>
                    <div className="text-xs text-gray-500">{stage.sector}</div>
                  </div>

                  <div className="text-xs text-gray-500 space-y-1">
                    <div>📍 {stage.commune}</div>
                    <div>
                      {stage.capacity_filled}/{stage.capacity_total} occupées ({tauxOccupation}%)
                    </div>
                    {stage.period && (
                      <div>📅 {stage.period}</div>
                    )}
                  </div>
                </div>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
