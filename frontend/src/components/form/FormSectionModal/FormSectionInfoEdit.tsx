'use client'

import type { FormSection } from '@/types/form'

interface FormSectionInfoEditProps {
    formSection: FormSection
    onChange: (updatedSection: FormSection) => void
}

const STATUS_OPTIONS = ['CREATED', 'PENDING', 'VALIDATED', 'REJECTED']

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
}

const CONTENT_LABELS: Record<string, string> = {
    key: 'Clé',
    description: 'Description',
}

function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString('fr-FR', {
        day: '2-digit', month: '2-digit', year: 'numeric',
        hour: '2-digit', minute: '2-digit'
    })
}

export default function FormSectionInfoEdit({ formSection, onChange }: FormSectionInfoEditProps) {

    const handleStatusChange = (value: string) => {
        onChange({ ...formSection, status: value })
    }

    const handleContentChange = (key: string, value: string) => {
        onChange({
            ...formSection,
            content: { ...formSection.content, [key]: value }
        })
    }

    return (
        <div className="px-6 pb-6 space-y-6">
            <div>
                <h4 className="text-lg font-medium text-gray-900 border-b pb-2 mb-4">
                    Section {SECTION_TYPE_LABELS[formSection.section_type] ?? formSection.section_type}
                </h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">

                    {/* Statut - éditable */}
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Statut</label>
                        <select
                            value={formSection.status}
                            onChange={(e) => handleStatusChange(e.target.value)}
                            className={`mt-1 px-2 py-1 text-xs font-semibold rounded-full border-0 ${STATUS_STYLES[formSection.status] ?? 'bg-gray-100 text-gray-800'}`}
                        >
                            {STATUS_OPTIONS.map((s) => (
                                <option key={s} value={s}>{s}</option>
                            ))}
                        </select>
                    </div>

                    {/* Utilisateur - non éditable */}
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Utilisateur</label>
                        <p className="text-sm text-gray-400">{formSection.user_id}</p>
                    </div>

                    {/* Dates - non éditables */}
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Créé le</label>
                        <p className="text-sm text-gray-400">{formatDate(formSection.created_at)}</p>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Mis à jour le</label>
                        <p className="text-sm text-gray-400">{formatDate(formSection.updated_at)}</p>
                    </div>

                    {/* Contenu mappé - éditable */}
                    {formSection.content && Object.keys(formSection.content).length > 0 && (
                        <div className="col-span-2">
                            <h5 className="text-md font-medium text-gray-900 pb-2 mb-4">Description</h5>
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                {Object.entries(formSection.content).map(([key, value]) => (
                                    <div key={key}>
                                        <label className="block text-sm font-medium text-gray-500">
                                            {CONTENT_LABELS[key] ?? key}
                                        </label>
                                        {typeof value === 'object' ? (
                                            <p className="text-sm text-gray-400">{JSON.stringify(value)}</p>
                                        ) : (
                                            <input
                                                type="text"
                                                value={String(value ?? '')}
                                                onChange={(e) => handleContentChange(key, e.target.value)}
                                                className="mt-1 block w-full text-sm border border-gray-300 rounded-md px-2 py-1 focus:outline-none focus:ring-1 focus:ring-blue-500"
                                            />
                                        )}
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
