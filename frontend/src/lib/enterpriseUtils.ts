// lib/enterpriseUtils.ts
import type { Tutor } from '@/types/tutor'

export interface TutorFormData {
    nom: string
    prenom: string
    email: string
    telephone?: string
    poste: string
}

export function validateTutorForm(data: TutorFormData): string[] {
    const errors: string[] = []

    if (!data.nom?.trim()) {
        errors.push('Le nom est requis')
    }

    if (!data.prenom?.trim()) {
        errors.push('Le prénom est requis')
    }

    if (!data.email?.trim()) {
        errors.push('L\'email est requis')
    } else if (!isValidEmail(data.email)) {
        errors.push('L\'email n\'est pas valide')
    }

    if (!data.poste?.trim()) {
        errors.push('Le poste est requis')
    }

    if (data.telephone && !isValidPhone(data.telephone)) {
        errors.push('Le numéro de téléphone n\'est pas valide')
    }

    return errors
}

export function isValidEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(email)
}

export function isValidPhone(phone: string): boolean {
    // Format français : 0X XX XX XX XX ou +33 X XX XX XX XX
    const phoneRegex = /^(?:(?:\+33|0)[1-9])(?:[ -]?\d{2}){4}$/
    return phoneRegex.test(phone.replace(/\s/g, ''))
}

export function formatTutorData(formData: TutorFormData): Omit<Tutor, 'id' | 'enterprise_id' | 'created_at' | 'updated_at'> {
    return {
        nom: formData.nom.trim(),
        prenom: formData.prenom.trim(),
        email: formData.email.trim().toLowerCase(),
        telephone: formData.telephone?.trim() || null,
        poste: formData.poste.trim(),
    }
}

export function getEmptyTutorForm(): TutorFormData {
    return {
        nom: '',
        prenom: '',
        email: '',
        telephone: '',
        poste: '',
    }
}

export function tutorToFormData(tutor: Tutor): TutorFormData {
    return {
        nom: tutor.nom,
        prenom: tutor.prenom,
        email: tutor.email,
        telephone: tutor.telephone || '',
        poste: tutor.poste,
    }
}

export function formatTutorDisplayName(tutor: Tutor): string {
    return `${tutor.prenom} ${tutor.nom}`
}

export function getTutorInitials(tutor: Tutor): string {
    return `${tutor.prenom?.[0]?.toUpperCase() || ''}${tutor.nom?.[0]?.toUpperCase() || ''}`
}
