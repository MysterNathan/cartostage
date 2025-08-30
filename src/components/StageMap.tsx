'use client'
import { useEffect, useRef } from 'react'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import type { Stage } from '@/types/stage'

// Configuration de l'icône personnalisée avec couleur selon filière
const createCustomIcon = (stage: Stage, isSelected = false) => {
  const filiereColor = stage.filiere === 'CCST' ? '#3B82F6' : '#10B981'
  const color = isSelected ? '#EF4444' : filiereColor
  
  return L.divIcon({
    html: `
      <div style="
        background-color: ${color};
        width: 26px;
        height: 26px;
        border-radius: 50%;
        border: 3px solid white;
        box-shadow: 0 2px 6px rgba(0,0,0,0.3);
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        font-size: 11px;
        font-weight: bold;
      ">
        ${stage.filiere === 'CCST' ? '⚡' : '❤️'}
      </div>
    `,
    className: 'custom-div-icon',
    iconSize: [26, 26],
    iconAnchor: [13, 26],
  })
}

interface StageMapProps {
  stages: Stage[]
  selectedStage?: Stage | null
  onStageClick?: (stage: Stage) => void
}

export default function StageMap({ stages, selectedStage, onStageClick }: StageMapProps) {
  const mapRef = useRef<L.Map | null>(null)
  const markersRef = useRef<{ [key: number]: L.Marker }>({})

  useEffect(() => {
    if (!mapRef.current) {
      mapRef.current = L.map('map', {
        center: [-21.1151, 55.5364],
        zoom: 10,
        scrollWheelZoom: true,
      })

      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© OpenStreetMap contributors'
      }).addTo(mapRef.current)
    }

    // Nettoyage des marqueurs existants  
    Object.values(markersRef.current).forEach(marker => {
      mapRef.current?.removeLayer(marker)
    })
    markersRef.current = {}

    // Ajout des nouveaux marqueurs
    stages.forEach((stage) => {
      const isSelected = selectedStage?.id === stage.id
      
      const marker = L.marker([stage.lat, stage.lng], {
        icon: createCustomIcon(stage, isSelected)
      }).addTo(mapRef.current!)

      // Popup enrichie avec tous les champs
      const tauxOccupation = stage.capacity_total > 0 
        ? Math.round((stage.capacity_filled / stage.capacity_total) * 100)
        : 0

      const popupContent = `
        <div style="min-width: 280px; font-family: ui-sans-serif, system-ui;">
          <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 8px;">
            <span style="
              background: ${stage.filiere === 'CCST' ? '#3B82F6' : '#10B981'}; 
              color: white; 
              padding: 2px 8px; 
              border-radius: 12px; 
              font-size: 11px; 
              font-weight: bold;
            ">${stage.filiere}</span>
          </div>
          
          <h3 style="margin: 0 0 8px 0; font-size: 16px; font-weight: bold; color: #111827;">
            ${stage.poste}
          </h3>
          
          <div style="margin-bottom: 8px;">
            <div style="font-weight: 600; color: #374151; margin-bottom: 2px;">
              ${stage.entreprise}
            </div>
            <div style="font-size: 12px; color: #6B7280;">
              ${stage.sector}
            </div>
          </div>

          <div style="font-size: 13px; color: #374151; margin-bottom: 8px;">
            📍 ${stage.commune}<br/>
            ${stage.adresse}
          </div>

          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 8px; margin-bottom: 8px;">
            <div style="background: #F3F4F6; padding: 6px; border-radius: 6px; text-align: center;">
              <div style="font-size: 11px; color: #6B7280;">Places libres</div>
              <div style="font-size: 16px; font-weight: bold; color: ${stage.placesDisponibles > 0 ? '#059669' : '#DC2626'};">
                ${stage.placesDisponibles}
              </div>
            </div>
            <div style="background: #F3F4F6; padding: 6px; border-radius: 6px; text-align: center;">
              <div style="font-size: 11px; color: #6B7280;">Taux occupation</div>
              <div style="font-size: 16px; font-weight: bold; color: #374151;">
                ${tauxOccupation}%
              </div>
            </div>
          </div>

          <div style="font-size: 12px; color: #6B7280;">
            <div><strong>Capacité:</strong> ${stage.capacity_filled}/${stage.capacity_total} occupées</div>
            ${stage.period ? `<div><strong>Périodes:</strong> ${stage.period}</div>` : ''}
          </div>
        </div>
      `

      marker.bindPopup(popupContent)

      // Gestion du clic
      if (onStageClick) {
        marker.on('click', () => {
          onStageClick(stage)
        })
      }

      markersRef.current[stage.id] = marker
    })

    // Auto-zoom sur les marqueurs visibles
    if (stages.length > 0) {
      const bounds = L.latLngBounds(stages.map(stage => [stage.lat, stage.lng]))
      if (bounds.isValid()) {
        mapRef.current.fitBounds(bounds.pad(0.1))
      }
    }

    // Ouvrir popup et centrer sur le stage sélectionné
    if (selectedStage && markersRef.current[selectedStage.id]) {
      markersRef.current[selectedStage.id].openPopup()
      mapRef.current.setView([selectedStage.lat, selectedStage.lng], 
        Math.max(mapRef.current.getZoom(), 13)
      )
    }

  }, [stages, selectedStage, onStageClick])

  return <div id="map" className="w-full h-full rounded-lg" />
}
