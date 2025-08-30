export interface Stage {
  id: number
  poste: string
  adresse: string
  lat: number
  lng: number
  placesDisponibles: number
  entreprise: string
  filiere: string
  sector: string
  commune: string
  capacity_total: number
  capacity_filled: number
  period?: string
}
