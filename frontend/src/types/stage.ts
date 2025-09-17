export interface Stage {
  id: number
  poste: string
  adresse: string
  lat: number
  lng: number
  places_disponibles: number
  enterprise: string
  filiere: string
  sector: string
  commune: string
  capacity_total: number
  capacity_filled: number
  period?: string
  parcours: "scolaire" | "apprentissage" | "mixte"
  famille_metiers: string
  niveau_scolaire: "2de" | "1re" | "Tle"
  created_at?: string // ISO datetime du backend
  updated_at?: string // ISO datetime du backend
}

export interface StagesData {
  stages: Stage[]
}
