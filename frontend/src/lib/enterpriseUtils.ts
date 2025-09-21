// lib/enterpriseUtils.ts
import type { Tutor } from '@/types/tutor'
import type { User } from '@/types/user'

// Types pour les tuteurs
export interface TutorFormData {
    nom: string
    prenom: string
    email: string
    telephone?: string
    poste: string
}

// Types pour les étudiants
export interface StudentFormData {
    prenom: string
    nom: string
    email: string
    telephone: string
    filiere: string
    etablissement: string
    niveau: string
    periode: string
    actif: string // "true" ou "false" pour le checkbox
}

// Validation des tuteurs
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

// Validation des étudiants
export function validateStudentForm(data: StudentFormData): string[] {
    const errors: string[] = []

    if (!data.prenom?.trim()) {
        errors.push('Le prénom est requis')
    }

    if (!data.nom?.trim()) {
        errors.push('Le nom est requis')
    }

    if (!data.email?.trim()) {
        errors.push('L\'email est requis')
    } else if (!isValidEmail(data.email)) {
        errors.push('L\'email n\'est pas valide')
    }

    if (!data.filiere?.trim()) {
        errors.push('La filière est requise')
    }

    if (!data.etablissement?.trim()) {
        errors.push('L\'établissement est requis')
    }

    if (data.telephone && !isValidPhone(data.telephone)) {
        errors.push('Le numéro de téléphone n\'est pas valide')
    }

    return errors
}

// Utilitaires de validation
export function isValidEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(email)
}

export function isValidPhone(phone: string): boolean {
    // Format français : 0X XX XX XX XX ou +33 X XX XX XX XX
    const phoneRegex = /^(?:(?:\+33|0)[1-9])(?:[ -]?\d{2}){4}$/
    return phoneRegex.test(phone.replace(/\s/g, ''))
}

// Formatage des données tuteurs
export function formatTutorData(formData: TutorFormData): Omit<Tutor, 'id' | 'enterprise_id' | 'created_at' | 'updated_at'> {
    return {
        nom: formData.nom.trim(),
        prenom: formData.prenom.trim(),
        email: formData.email.trim().toLowerCase(),
        telephone: formData.telephone?.trim() || null,
        poste: formData.poste.trim(),
    }
}

// Formatage des données étudiants
export function formatStudentData(formData: StudentFormData) {
    return {
        first_name: formData.prenom.trim(),
        last_name: formData.nom.trim(),
        email: formData.email.trim().toLowerCase(),
        is_active: formData.actif === 'true',
        profile: {
            phone: formData.telephone.trim() || null,
            filiere: formData.filiere.trim(),
            etablissement: formData.etablissement.trim(),
            niveau: formData.niveau || null,
            periode: formData.periode.trim() || null,
            is_active: formData.actif === 'true'
        }
    }
}

// Formulaires vides
export function getEmptyTutorForm(): TutorFormData {
    return {
        nom: '',
        prenom: '',
        email: '',
        telephone: '',
        poste: '',
    }
}

export function getEmptyStudentForm(): StudentFormData {
    return {
        prenom: '',
        nom: '',
        email: '',
        telephone: '',
        filiere: '',
        etablissement: '',
        niveau: '',
        periode: '',
        actif: 'true'
    }
}

// Conversion vers formulaire
export function tutorToFormData(tutor: Tutor): TutorFormData {
    return {
        nom: tutor.nom,
        prenom: tutor.prenom,
        email: tutor.email,
        telephone: tutor.telephone || '',
        poste: tutor.poste,
    }
}

export function studentToFormData(student: User): StudentFormData {
    return {
        prenom: student.first_name || '',
        nom: student.last_name || '',
        email: student.email || '',
        telephone: student.profile?.phone || '',
        filiere: student.profile?.filiere || '',
        etablissement: student.profile?.etablissement || '',
        niveau: student.profile?.niveau || '',
        periode: student.profile?.periode || '',
        actif: (student.is_active && (student.profile?.is_active !== false)).toString()
    }
}

// Utilitaires d'affichage pour les tuteurs
export function formatTutorDisplayName(tutor: Tutor): string {
    return `${tutor.prenom} ${tutor.nom}`
}

export function getTutorInitials(tutor: Tutor): string {
    return `${tutor.prenom?.[0]?.toUpperCase() || ''}${tutor.nom?.[0]?.toUpperCase() || ''}`
}

// Utilitaires d'affichage pour les étudiants
export function formatStudentDisplayName(student: User): string {
    return `${student.first_name || ''} ${student.last_name || ''}`.trim()
}

export function getStudentInitials(student: User): string {
    return `${student.first_name?.[0]?.toUpperCase() || ''}${student.last_name?.[0]?.toUpperCase() || ''}`
}

// Helpers pour les filières et établissements
export function formatFiliere(filiere: string): string {
    return filiere || 'Non spécifiée'
}

export function formatEtablissement(etablissement: string): string {
    return etablissement || 'Non spécifié'
}

export function formatNiveau(niveau: string): string {
    if (!niveau) return ''
    return niveau.charAt(0).toUpperCase() + niveau.slice(1)
}
