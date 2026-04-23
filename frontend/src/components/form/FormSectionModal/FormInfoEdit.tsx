'use client'

import type { Form } from '@/types/form'

interface FormInfoEditProps {
    form: Form
    onChange: (updatedForm: Form) => void
}

const CONTENT_LABELS: Record<string, string> = {
    key: 'Clé',
    description: 'Description',
}

const STATUS_OPTIONS = ['CREATED', 'PENDING', 'VALIDATED', 'REJECTED']

const STATUS_STYLES: Record<string, string> = {
    CREATED: 'bg-blue-100 text-blue-800',
    PENDING: 'bg-yellow-100 text-yellow-800',
    VALIDATED: 'bg-green-100 text-green-800',
    REJECTED: 'bg-red-100 text-red-800',
}

function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString('fr-FR', {
        day: '2-digit', month: '2-digit', year: 'numeric',
        hour: '2-digit', minute: '2-digit'
    })
}

export default function FormInfoEdit({ form, onChange }: FormInfoEditProps) {

    const handleStatusChange = (value: string) => {
        onChange({ ...form, status: value })
    }

    const handleContentChange = (key: string, value: string) => {
        onChange({
            ...form,
            content: { ...form.content, [key]: value }
        })
    }

    return (
        <div className="p-6 space-y-6">

            <div>
                <h4 className="text-lg font-medium text-gray-900 border-b pb-2 mb-4">Général</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">

                    {/* Statut - éditable */}
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Statut</label>
                        <select
                            value={form.status}
                            onChange={(e) => handleStatusChange(e.target.value)}
                            className={`mt-1 px-2 py-1 text-xs font-semibold rounded-full border-0 ${STATUS_STYLES[form.status] ?? 'bg-gray-100 text-gray-800'}`}
                        >
                            {STATUS_OPTIONS.map((s) => (
                                <option key={s} value={s}>{s}</option>
                            ))}
                        </select>
                    </div>

                    {/* Dates - non éditables */}
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Créé le</label>
                        <p className="text-sm text-gray-400">{formatDate(form.created_at)}</p>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Mis à jour le</label>
                        <p className="text-sm text-gray-400">{formatDate(form.updated_at)}</p>
                    </div>

                </div>
            </div>

            {/* Contenu mappé - éditable */}
            {form.content && Object.keys(form.content).length > 0 && (
                <div>
                    <h5 className="text-md font-medium text-gray-900 pb-2 mb-4">Description</h5>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {Object.entries(form.content).map(([key, value]) => (
                            <div key={key}>
                                <label className="block text-sm font-medium text-gray-500">
                                    {CONTENT_LABELS[key] ?? key}
                                </label>
                                {typeof value === 'object' ? (
                                    // Les objets imbriqués restent en lecture seule
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
    )
}
