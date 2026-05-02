// components/form/FormInfo.tsx
'use client'

import type { Form } from '@/types/form'

interface FormInfoProps {
    form: Form
}

const CONTENT_LABELS: Record<string, string> = {
    key: 'Clé',
    description: 'Description',
    // Ajoute tes clés ici
}

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

export default function FormInfo({ form }: FormInfoProps) {
    return (
        <div className="p-6 space-y-6">

            {/* Informations du formulaire */}
            <div>
                <h4 className="text-lg font-medium text-gray-900 border-b pb-2 mb-4">Général</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Statut</label>
                        <span className={`mt-1 px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${STATUS_STYLES[form.status] ?? 'bg-gray-100 text-gray-800'}`}>
                            {form.status}
                        </span>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Créé le</label>
                        <p className="text-sm text-gray-900">{formatDate(form.created_at)}</p>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Mis à jour le</label>
                        <p className="text-sm text-gray-900">{formatDate(form.updated_at)}</p>
                    </div>
                </div>
            </div>

            {/* Contenu mappé */}
            {form.content && Object.keys(form.content).length > 0 && (
                <div>
                    <h5 className="text-md font-medium text-gray-900 pb-2 mb-4">Description</h5>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {Object.entries(form.content).map(([key, value]) => (
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
    )
}
