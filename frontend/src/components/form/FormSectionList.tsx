'use client'

import { useState } from 'react'
import type { FormSection } from '@/types/formSection'
import FormSectionModal from '@/components/form/FormSectionModal'

interface FormSectionsListProps {
    sections: FormSection[]
}

const STATUS_STYLES: Record<string, string> = {
    completed: 'bg-green-100 text-green-800',
    pending:   'bg-yellow-100 text-yellow-800',
    rejected:  'bg-red-100 text-red-800',
    draft:     'bg-gray-100 text-gray-800',
}

const STATUS_LABELS: Record<string, string> = {
    completed: 'Complété',
    pending:   'En attente',
    rejected:  'Rejeté',
    draft:     'Brouillon',
}

// Helpers
const getStatusStyle = (status: string): string =>
    STATUS_STYLES[status] ?? 'bg-gray-100 text-gray-800'

const getStatusLabel = (status: string): string =>
    STATUS_LABELS[status] ?? status

const formatDate = (dateStr: string): string =>
    new Date(dateStr).toLocaleDateString('fr-FR', {
        day:   '2-digit',
        month: '2-digit',
        year:  'numeric',
    })

export default function FormSectionsList({ sections }: FormSectionsListProps) {
    const [selectedSection, setSelectedSection] = useState<FormSection | null>(null)

    return (
        <>
            <div className="bg-white rounded-lg shadow overflow-hidden">
                <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
                    <h3 className="text-lg font-medium text-gray-900">
                        Sections ({sections?.length ?? 0})
                    </h3>
                </div>

                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                        <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Section
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Formulaire
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Statut
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Dates
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Actions
                            </th>
                        </tr>
                        </thead>

                        <tbody className="bg-white divide-y divide-gray-200">
                        {sections.map((section) => (
                            <tr key={section.id} className="hover:bg-gray-50">

                                {/* Section type */}
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="flex items-center">
                                        <div className="flex-shrink-0 h-10 w-10">
                                            <div className="h-10 w-10 rounded-full bg-blue-500 flex items-center justify-center">
                                                    <span className="text-sm font-medium text-white">
                                                        {section.section_type[0]?.toUpperCase() ?? '?'}
                                                    </span>
                                            </div>
                                        </div>
                                        <div className="ml-4">
                                            <div className="text-sm font-medium text-gray-900">
                                                {section.section_type}
                                            </div>
                                            <div className="text-sm text-gray-500">
                                                #{section.id}
                                            </div>
                                        </div>
                                    </div>
                                </td>

                                {/* Form ID */}
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="text-sm text-gray-900">
                                        Formulaire #{section.form_id}
                                    </div>
                                    <div className="text-sm text-gray-500">
                                        Utilisateur #{section.user_id}
                                    </div>
                                </td>

                                {/* Statut */}
                                <td className="px-6 py-4 whitespace-nowrap">
                                        <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusStyle(section.status)}`}>
                                            {getStatusLabel(section.status)}
                                        </span>
                                    {section.completed_at && (
                                        <div className="text-xs text-gray-500 mt-1">
                                            Complété le {formatDate(section.completed_at)}
                                        </div>
                                    )}
                                </td>

                                {/* Dates */}
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="text-sm text-gray-900">
                                        Créé le {formatDate(section.created_at)}
                                    </div>
                                    <div className="text-sm text-gray-500">
                                        Modifié le {formatDate(section.updated_at)}
                                    </div>
                                </td>

                                {/* Actions */}
                                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                                    <div className="flex items-center space-x-3">
                                        <button
                                            onClick={() => setSelectedSection(section)}
                                            className="text-blue-600 hover:text-blue-900 font-medium flex items-center"
                                        >
                                            <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                                            </svg>
                                            Voir le contenu
                                        </button>
                                    </div>
                                </td>

                            </tr>
                        ))}
                        </tbody>
                    </table>
                </div>
            </div>

            <FormSectionModal
                section={selectedSection}
                onClose={() => setSelectedSection(null)}
            />
        </>
    )
}
