// components/form/FormSectionInfo.tsx
'use client'

import type { FormSection } from '@/types/form'

interface FormSectionInfoProps {
    formSection: FormSection
}

const STATUS_STYLES: Record<string, string> = {
    CREATED: 'bg-blue-100 text-blue-800',
    PENDING: 'bg-yellow-100 text-yellow-800',
    VALIDATED: 'bg-green-100 text-green-800',
    REJECTED: 'bg-red-100 text-red-800',
}

const SECTION_TYPE_LABELS: Record<string, string> = {
    STUDENT: 'Étudiant',
    TUTOR: 'Tuteur',
    SCHOOL: 'École',
    // Ajoute tes types ici
}

const CONTENT_LABELS: Record<string, string> = {
    key: 'Clé',
    description: 'Description',
    // Ajoute tes clés ici
}

function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString('fr-FR', {
        day: '2-digit', month: '2-digit', year: 'numeric',
        hour: '2-digit', minute: '2-digit'
    })
}

export default function FormSectionInfo({ formSection }: FormSectionInfoProps) {
    return (
        <div className="px-6 pb-6 space-y-6">

            <div>
                <h4 className="text-lg font-medium text-gray-900 border-b pb-2 mb-4">Section {SECTION_TYPE_LABELS[formSection.section_type] ?? formSection.section_type}</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Statut</label>
                        <span className={`mt-1 px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${STATUS_STYLES[formSection.status] ?? 'bg-gray-100 text-gray-800'}`}>
                            {formSection.status}
                        </span>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Utilisateur</label>
                        <p className="text-sm text-gray-900">{formSection.user_id}</p>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Créé le</label>
                        <p className="text-sm text-gray-900">{formatDate(formSection.created_at)}</p>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Mis à jour le</label>
                        <p className="text-sm text-gray-900">{formatDate(formSection.updated_at)}</p>

                    </div>
                    {/* Contenu mappé */}
                    {formSection.content && Object.keys(formSection.content).length > 0 && (
                        <div>
                            <h5 className="text-md font-medium text-gray-900 pb-2 mb-4">Description</h5>
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                {Object.entries(formSection.content).map(([key, value]) => (
                                    <div key={key}>
                                        <label className="block text-sm font-medium text-gray-500">
                                            {CONTENT_LABELS[key] ?? key}
                                        </label>
                                        <p className="text-sm text-gray-900">
                                            {typeof value === 'object'
                                                ? JSON.stringify(value)
                                                : String(value ?? '—')}
                                        </p>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </div>
            </div>

        </div>
    )
}
