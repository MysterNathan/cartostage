export interface Enterprise {
    id: number
    nom: string
    adresse: string
    secteur: string
    taille: string
    siret: string
    email: string
    telephone: string
    site_web: string
    description: string
    logo_url: string
    created_at: string
    updated_at: string
}

export interface EnterpriseWithStats extends Enterprise {
    total_tutors: number
    active_stages: number
    total_students: number
}
