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
    const [isFullscreen, setIsFullscreen] = useState(false);


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
        <div className={`fixed inset-0 bg-white/30 backdrop-blur-sm flex z-50 ${
            isFullscreen ? 'items-start p-0' : 'items-center justify-center p-4'
        }`}>
            <div className={`bg-white shadow-xl border border-gray-200 w-full overflow-y-auto transition-all duration-200 ${
                isFullscreen ? 'max-w-full h-full rounded-none' : 'max-w-4xl max-h-[90vh] rounded-lg'
            }`}>

                {/* Header — HORS de la boucle */}
                <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center sticky top-0 bg-white z-10">
                    <div className="flex items-center gap-3">
                        <h3 className="text-lg font-medium text-gray-900">
                            Fiches de stage
                        </h3>
                        {isEditing && (
                            <span className="text-xs font-medium text-blue-600 bg-blue-50 px-2 py-1 rounded-full">
                Mode édition
            </span>
                        )}
                    </div>
                    <div className="flex items-center gap-2">
                        <button onClick={() => setIsFullscreen(!isFullscreen)} className="text-gray-400 hover:text-gray-600">
                            {isFullscreen ? (
                                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 9L4 4m0 0h5m-5 0v5M15 9l5-5m0 0h-5m5 0v5M9 15l-5 5m0 0h5m-5 0v-5M15 15l5 5m0 0h-5m5 0v-5" />
                                </svg>
                            ) : (
                                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 8V4m0 0h4M4 4l5 5M20 8V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5M20 16v4m0 0h-4m4 0l-5-5" />
                                </svg>
                            )}
                        </button>
                        <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
                            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                </div>


                {/* Contenu — boucle sur les fiches */}
                <div className="bg-gray-200 p-4 space-y-4">
                    {editedData.map(({ form, form_section }) => (
                        <div key={form.id} className="bg-white border border-gray-200 rounded-lg p-6">
                            <h4 className="text-sm font-semibold text-gray-700 mb-4">
                                Fiche n°{form.id}
                            </h4>

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
                </div>

                {/* Footer — HORS de la boucle */}
                <div className="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-between sticky bottom-0">
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
