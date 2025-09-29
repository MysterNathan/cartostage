export interface StageOffer {
    id: number
    position: string
    address: string
    lat: number
    lng: number
    enterprise: string
    sector: string
    capacity_total: number
    capacity_filled: number
    period: string
    course: string
    job_family: string
    scolar_level: string
    created_at: string // ISO datetime
    updated_at: string // ISO datetime
}

// Type dérivé pour l'affichage sur la carte (si nécessaire)
export interface StageOfferForMap {
    id: number
    poste: string // position
    adresse: string // address
    lat: number
    lng: number
    places_disponibles: number // capacity_total - capacity_filled
    enterprise: string
    filiere: string // course
    sector: string
    commune: string // extrait de address
    capacity_total: number
    capacity_filled: number
    period: string
    famille_metiers: string // job_family
    niveau_scolaire: string // scolar_level
    created_at: string
    updated_at: string
}
