export interface Filiere {
    id: number
    code: string
    label: string
    color: string
}

export interface FilieresResponse {
    filieres: Filiere[]
}

export interface CreateFiliereData {
    code: string
    label: string
    color: string
}

export interface UpdateFiliereData {
    code?: string
    label?: string
    color?: string
}
