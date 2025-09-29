import {StageOffer} from "@/types/stage_offer";

export interface Stage {
  id: number
  stage_offer_id: number
  student_id: number
  teacher_id?: number
  tutor_id?: number
  establishment_id?: number
  content_id?: number
  status: string
  start_date: string // ISO datetime
  end_date: string // ISO datetime
  created_at: string // ISO datetime
  updated_at: string // ISO datetime
}

export interface StageWithDetails {
  id: number
  stage_offer_id: number
  student_id: number
  teacher_id?: number
  tutor_id?: number
  establishment_id?: number
  content_id?: number
  status: string
  start_date: string
  end_date: string
  created_at: string
  updated_at: string
  stage_offer?: StageOffer
  student?: UserPublic
  teacher?: UserPublic
  tutor?: UserPublic
  establishment?: Establishment
  content?: Content
}

export interface CreateStageRequest {
  stage_offer_id: number
  student_id: number
  teacher_id?: number
  tutor_id?: number
  establishment_id?: number
  content_id?: number
  status: string
  start_date: string
  end_date: string
}

export interface UpdateStageRequest {
  teacher_id?: number
  tutor_id?: number
  content_id?: number
  status?: string
  start_date?: string
  end_date?: string
}

export interface Content {
  id: number
  content?: string
}

// Types utilitaires
export interface UserPublic {
  id: number
  // Ajouter les autres champs selon votre modèle User
}

export interface Establishment {
  id: number
  // Ajouter les autres champs selon votre modèle Establishment
}
