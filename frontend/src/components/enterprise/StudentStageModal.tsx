// components/enterprise/StudentStageModal.tsx
'use client'

import { formatStudentDisplayName } from '@/lib/enterpriseUtils'
import type { User } from '@/types/user'

interface StudentStageModalProps {
    student: User | null
    onClose: () => void
}

export default function StudentStageModal({ student, onClose }: StudentStageModalProps) {
    if (!student) return null

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
                <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
                    <h3 className="text-lg font-medium text-gray-900">
                        Fiche de stage - {formatStudentDisplayName(student)}
                    </h3>
                    <button
                        onClick={onClose}
                        className="text-gray-400 hover:text-gray-600"
                    >
                        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                <div className="p-6">
                    <div className="bg-yellow-50 border border-yellow-200 rounded-md p-4 mb-6">
                        <div className="flex">
                            <svg className="h-5 w-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
                                <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                            </svg>
                            <div className="ml-3">
                                <h3 className="text-sm font-medium text-yellow-800">
                                    Fonctionnalité en développement
                                </h3>
                                <p className="mt-1 text-sm text-yellow-700">
                                    La gestion des stages est actuellement en cours de développement.
                                </p>
                            </div>
                        </div>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="space-y-4">
                            <h4 className="text-md font-medium text-gray-900 border-b pb-2">Informations de l'élève</h4>
                            <div>
                                <label className="block text-sm font-medium text-gray-500">Nom complet</label>
                                <p className="text-sm text-gray-900">{formatStudentDisplayName(student)}</p>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-500">Filière</label>
                                <p className="text-sm text-gray-900">{student.profile?.filiere || 'Non spécifiée'}</p>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-500">Établissement</label>
                                <p className="text-sm text-gray-900">{student.profile?.etablissement || 'Non spécifié'}</p>
                            </div>
                        </div>

                        <div className="space-y-4">
                            <h4 className="text-md font-medium text-gray-900 border-b pb-2">Stage (à venir)</h4>
                            <div className="bg-gray-50 rounded-lg p-4 text-center">
                                <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                                </svg>
                                <p className="mt-2 text-sm text-gray-500">
                                    Les informations de stage seront disponibles prochainement
                                </p>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="px-6 py-4 bg-gray-50 flex justify-end">
                    <button
                        onClick={onClose}
                        className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                    >
                        Fermer
                    </button>
                </div>
            </div>
        </div>
    )
}
