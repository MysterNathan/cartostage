// components/enterprise/StudentForm.tsx
'use client'

import { useState } from 'react'
import { validateStudentForm, getEmptyStudentForm, studentToFormData, formatStudentData, type StudentFormData } from '@/lib/enterpriseUtils'
import type { User } from '@/types/user'

interface StudentFormProps {
    student?: User | null
    onSubmit: (data: ReturnType<typeof formatStudentData>) => Promise<void>
    onCancel: () => void
    loading?: boolean
}

export default function StudentForm({ student, onSubmit, onCancel, loading }: StudentFormProps) {
    const [formData, setFormData] = useState<StudentFormData>(
        student ? studentToFormData(student) : getEmptyStudentForm()
    )
    const [errors, setErrors] = useState<string[]>([])
    const [isSubmitting, setIsSubmitting] = useState(false)

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        const validationErrors = validateStudentForm(formData)
        if (validationErrors.length > 0) {
            setErrors(validationErrors)
            return
        }

        setErrors([])
        setIsSubmitting(true)

        try {
            await onSubmit(formatStudentData(formData))
        } catch (error) {
            setErrors([error instanceof Error ? error.message : 'Erreur lors de la sauvegarde'])
        } finally {
            setIsSubmitting(false)
        }
    }

    const handleChange = (field: keyof StudentFormData, value: string) => {
        setFormData(prev => ({ ...prev, [field]: value }))
        if (errors.length > 0) {
            setErrors([])
        }
    }

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-lg shadow-xl max-w-md w-full max-h-[90vh] overflow-y-auto">
                <div className="px-6 py-4 border-b border-gray-200">
                    <h3 className="text-lg font-medium text-gray-900">
                        {student ? 'Modifier l\'élève' : 'Ajouter un élève'}
                    </h3>
                </div>

                <form onSubmit={handleSubmit} className="p-6 space-y-4">
                    {errors.length > 0 && (
                        <div className="bg-red-50 border border-red-200 rounded-md p-4">
                            <div className="flex">
                                <svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                                </svg>
                                <div className="ml-3">
                                    <h3 className="text-sm font-medium text-red-800">
                                        Erreurs de validation
                                    </h3>
                                    <div className="mt-2 text-sm text-red-700">
                                        <ul className="list-disc list-inside space-y-1">
                                            {errors.map((error, index) => (
                                                <li key={index}>{error}</li>
                                            ))}
                                        </ul>
                                    </div>
                                </div>
                            </div>
                        </div>
                    )}

                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                Prénom *
                            </label>
                            <input
                                type="text"
                                value={formData.prenom}
                                onChange={(e) => handleChange('prenom', e.target.value)}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-green-500 focus:border-green-500"
                                disabled={isSubmitting}
                                required
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                Nom *
                            </label>
                            <input
                                type="text"
                                value={formData.nom}
                                onChange={(e) => handleChange('nom', e.target.value)}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-green-500 focus:border-green-500"
                                disabled={isSubmitting}
                                required
                            />
                        </div>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Email *
                        </label>
                        <input
                            type="email"
                            value={formData.email}
                            onChange={(e) => handleChange('email', e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-green-500 focus:border-green-500"
                            disabled={isSubmitting}
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Téléphone
                        </label>
                        <input
                            type="tel"
                            value={formData.telephone}
                            onChange={(e) => handleChange('telephone', e.target.value)}
                            placeholder="0X XX XX XX XX"
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-green-500 focus:border-green-500"
                            disabled={isSubmitting}
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Filière *
                        </label>
                        <input
                            type="text"
                            value={formData.filiere}
                            onChange={(e) => handleChange('filiere', e.target.value)}
                            placeholder="Ex: BTS SIO, DUT Informatique, Licence Pro..."
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-green-500 focus:border-green-500"
                            disabled={isSubmitting}
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Établissement *
                        </label>
                        <input
                            type="text"
                            value={formData.etablissement}
                            onChange={(e) => handleChange('etablissement', e.target.value)}
                            placeholder="Ex: Lycée Jean Moulin, IUT de Bordeaux..."
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-green-500 focus:border-green-500"
                            disabled={isSubmitting}
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Niveau d'étude
                        </label>
                        <select
                            value={formData.niveau}
                            onChange={(e) => handleChange('niveau', e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-green-500 focus:border-green-500"
                            disabled={isSubmitting}
                        >
                            <option value="">Sélectionner un niveau</option>
                            <option value="bac">Bac</option>
                            <option value="bac+1">Bac+1</option>
                            <option value="bac+2">Bac+2</option>
                            <option value="bac+3">Bac+3</option>
                            <option value="bac+4">Bac+4</option>
                            <option value="bac+5">Bac+5</option>
                        </select>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Période de stage
                        </label>
                        <input
                            type="text"
                            value={formData.periode}
                            onChange={(e) => handleChange('periode', e.target.value)}
                            placeholder="Ex: Du 01/03/2024 au 30/06/2024"
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-green-500 focus:border-green-500"
                            disabled={isSubmitting}
                        />
                    </div>

                    <div className="flex items-center">
                        <input
                            type="checkbox"
                            id="actif"
                            checked={formData.actif}
                            onChange={(e) => handleChange('actif', e.target.checked.toString())}
                            className="h-4 w-4 text-green-600 focus:ring-green-500 border-gray-300 rounded"
                            disabled={isSubmitting}
                        />
                        <label htmlFor="actif" className="ml-2 block text-sm text-gray-900">
                            Élève actif
                        </label>
                    </div>

                    <div className="flex justify-end space-x-3 pt-4">
                        <button
                            type="button"
                            onClick={onCancel}
                            disabled={isSubmitting}
                            className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"
                        >
                            Annuler
                        </button>
                        <button
                            type="submit"
                            disabled={isSubmitting}
                            className="px-4 py-2 text-sm font-medium text-white bg-green-600 border border-transparent rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
                        >
                            {isSubmitting && (
                                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                            )}
                            {student ? 'Modifier' : 'Ajouter'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    )
}
