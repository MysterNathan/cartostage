'use client'

import { useState } from 'react'
import type { FormResponse } from '@/types/form'
import FormInfo from "@/components/form/FormSectionModal/FormInfo"
import FormSectionInfo from "@/components/form/FormSectionModal/FormSectionInfo"
import FormInfoEdit from "@/components/form/FormSectionModal/FormInfoEdit"
import FormSectionInfoEdit from "@/components/form/FormSectionModal/FormSectionInfoEdit"
import {updateForm} from "@/lib/api/stageApi";

interface FormSectionModalProps {
    formResponse: FormResponse
    onClose: () => void
    onSave?: (updatedData: FormResponse) => Promise<void>
}

export default function FormSectionModal({ formResponse, onClose }: FormSectionModalProps) {
    const [isEditing, setIsEditing] = useState(false)
    const [isSaving, setIsSaving] = useState(false)
    const [editedData, setEditedData] = useState<FormResponse>(formResponse)

    if (formResponse.error) {
        return (
            <div className="fixed inset-0 bg-white/30 backdrop-blur-sm flex items-center justify-center p-4 z-50">
                <div className="bg-white rounded-lg shadow-xl border border-gray-200 max-w-md w-full p-6 text-center">
                    <h3 className="text-lg font-medium text-gray-900 mb-2">
                        Aucune donnée disponible
                    </h3>
                    <p className="text-sm text-gray-600 mb-4">
                        Les informations de cette fiche ne sont pas disponibles.
                    </p>
                    <button
                        onClick={onClose}
                        className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                    >
                        Fermer
                    </button>
                </div>
            </div>
        )
    }

    const handleCancel = () => {
        setEditedData(formResponse)
        setIsEditing(false)
    }

    const handleSave = async () => {
        setIsSaving(true)
        try {
            await updateForm(editedData[0])
            onClose()
        } catch (error) {
            console.error('Erreur lors de la sauvegarde:', error)
        } finally {
            setIsSaving(false)
        }
    }

    return (
        <div className="fixed inset-0 bg-white/30 backdrop-blur-sm flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-lg shadow-xl border border-gray-200 max-w-2xl w-full max-h-[90vh] overflow-y-auto">

                {editedData.map(({ form, form_section }) => (
                    <div key={form.id}>
                        {/* Header */}
                        <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
                            <div className="flex items-center gap-3">
                                <h3 className="text-lg font-medium text-gray-900">
                                    Fiche de stage n°{form.id}
                                </h3>
                                {isEditing && (
                                    <span className="text-xs font-medium text-blue-600 bg-blue-50 px-2 py-1 rounded-full">
                                        Mode édition
                                    </span>
                                )}
                            </div>
                            <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
                                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                </svg>
                            </button>
                        </div>

                        {/* Contenu : lecture ou édition */}
                        {isEditing ? (
                            <>
                                <FormInfoEdit
                                    form={form}
                                    onChange={(updatedForm) =>
                                        setEditedData(prev =>
                                            prev.map(item =>
                                                item.form.id === form.id
                                                    ? { ...item, form: updatedForm }
                                                    : item
                                            )
                                        )
                                    }
                                />
                                {form_section.map((section) => (
                                    <FormSectionInfoEdit
                                        key={section.id}
                                        formSection={section}
                                        onChange={(updatedSection) =>
                                            setEditedData(prev =>
                                                prev.map(item =>
                                                    item.form.id === form.id
                                                        ? {
                                                            ...item,
                                                            form_section: item.form_section.map(s =>
                                                                s.id === section.id ? updatedSection : s
                                                            )
                                                        }
                                                        : item
                                                )
                                            )
                                        }
                                    />
                                ))}
                            </>
                        ) : (
                            <>
                                <FormInfo form={form} />
                                {form_section.map((section) => (
                                    <FormSectionInfo key={section.id} formSection={section} />
                                ))}
                            </>
                        )}
                    </div>
                ))}

                {/* Footer */}
                <div className="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-between">
                    {/* Bouton Modifier (visible en mode lecture) */}
                    {!isEditing ? (
                        <button
                            onClick={() => setIsEditing(true)}
                            className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                        >
                            Modifier
                        </button>
                    ) : (
                        <div className="flex gap-2">
                            <button
                                onClick={handleCancel}
                                disabled={isSaving}
                                className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"
                            >
                                Annuler
                            </button>
                            <button
                                onClick={handleSave}
                                disabled={isSaving}
                                className="px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700 disabled:opacity-50"
                            >
                                {isSaving ? 'Sauvegarde...' : 'Sauvegarder'}
                            </button>
                        </div>
                    )}

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
