'use client'
import { useState } from 'react'
import type { Stage } from '@/types/stage'

interface StageFormProps {
  stage?: Stage
  onSubmit: (stage: Omit<Stage, 'id'>) => void
  onCancel: () => void
}

export default function StageForm({ stage, onSubmit, onCancel }: StageFormProps) {
  const [formData, setFormData] = useState({
    poste: stage?.poste || '',
    entreprise: stage?.entreprise || '',
    adresse: stage?.adresse || '',
    lat: stage?.lat || 0,
    lng: stage?.lng || 0,
    placesDisponibles: stage?.placesDisponibles || 1,
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(formData)
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: type === 'number' ? Number(value) : value
    }))
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h3 className="text-lg font-semibold mb-4">
        {stage ? 'Modifier le stage' : 'Nouveau stage'}
      </h3>
      
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">Poste</label>
          <input
            type="text"
            name="poste"
            value={formData.poste}
            onChange={handleChange}
            required
            className="w-full border border-gray-300 rounded px-3 py-2"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Entreprise</label>
          <input
            type="text"
            name="entreprise"
            value={formData.entreprise}
            onChange={handleChange}
            required
            className="w-full border border-gray-300 rounded px-3 py-2"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Adresse</label>
          <input
            type="text"
            name="adresse"
            value={formData.adresse}
            onChange={handleChange}
            required
            className="w-full border border-gray-300 rounded px-3 py-2"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1">Latitude</label>
            <input
              type="number"
              name="lat"
              value={formData.lat}
              onChange={handleChange}
              step="0.0001"
              required
              className="w-full border border-gray-300 rounded px-3 py-2"
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Longitude</label>
            <input
              type="number"
              name="lng"
              value={formData.lng}
              onChange={handleChange}
              step="0.0001"
              required
              className="w-full border border-gray-300 rounded px-3 py-2"
            />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Places disponibles</label>
          <input
            type="number"
            name="placesDisponibles"
            value={formData.placesDisponibles}
            onChange={handleChange}
            min="0"
            required
            className="w-full border border-gray-300 rounded px-3 py-2"
          />
        </div>

        <div className="flex gap-2 pt-4">
          <button
            type="submit"
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            {stage ? 'Modifier' : 'Ajouter'}
          </button>
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
          >
            Annuler
          </button>
        </div>
      </form>
    </div>
  )
}
