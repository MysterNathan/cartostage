'use client'

import type { FormResponse } from '@/types/form'
import FormInfo from "@/components/form/FormSectionModal/FormInfo";
import FormSectionInfo from "@/components/form/FormSectionModal/FormSectionInfo";

interface FormSectionModalProps {
    formResponse: FormResponse
    onClose: () => void
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

export default function FormSectionModal({ formResponse, onClose }: FormSectionModalProps) {
    const { form, form_section } = formResponse

    return (
        <div className="fixed inset-0 bg-white/30 backdrop-blur-sm flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-lg shadow-xl border border-gray-200 max-w-2xl w-full max-h-[90vh] overflow-y-auto">

                {/* Header */}
                <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
                    <h3 className="text-lg font-medium text-gray-900">
                        Fiche de stage n°{form.id}
                    </h3>
                    <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
                        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                <FormInfo form={form} />
                <FormSectionInfo formSection={form_section} />


                {/* Footer */}
                <div className="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-end">
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
