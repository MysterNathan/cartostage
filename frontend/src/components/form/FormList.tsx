// components/enterprise-service/FormsList.tsx
'use client'

import { useState, useEffect } from 'react'
import type {FormResponse} from '@/types/form'

interface FormsListProps {
    initialForms?: FormResponses
    autoFetch?: boolean
}

export default function FormsList({ initialForms, autoFetch = true }: FormsListProps) {
    const [forms, setForms] = useState<FormResponse[]>(initialForms?.form_responses || [])
    const [loading, setLoading] = useState(autoFetch && !initialForms)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        if (autoFetch && !initialForms) {
            loadForms()
        }
    }, [autoFetch, initialForms])

    const loadForms = async () => {
        try {
            setLoading(true)
            setError(null)
            const data = await getForms()
            setForms(data.form_responses || [])
        } catch (error) {
            console.error('Erreur lors du chargement des formulaires:', error)
            setError(error instanceof Error ? error.message : 'Erreur de chargement')
        } finally {
            setLoading(false)
        }
    }


    const getStatusStyle = (status: string) => {
        switch (status) {
            case 'approved': return 'bg-green-100 text-green-800'
            case 'pending': return 'bg-yellow-100 text-yellow-800'
            case 'rejected': return 'bg-red-100 text-red-800'
            default: return 'bg-gray-100 text-gray-800'
        }
    }

    const formatDate = (date: string) => {
        return new Date(date).toLocaleDateString('fr-FR')
    }

    if (error && !loading) {
        return (
            <div className="bg-white rounded-lg shadow">
                <div className="p-6">
                    <div className="text-center">
                        <svg className="mx-auto h-12 w-12 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <h3 className="mt-2 text-sm font-medium text-gray-900">Erreur de chargement</h3>
                        <p className="mt-1 text-sm text-gray-500">{error}</p>
                        <button
                            onClick={loadForms}
                            className="mt-4 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                        >
                            Réessayer
                        </button>
                    </div>
                </div>
            </div>
        )
    }

    if (loading) {
        return (
            <div className="bg-white rounded-lg shadow">
                <div className="p-6">
                    <div className="text-center">
                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-4"></div>
                        <p className="text-gray-600">Chargement des formulaires...</p>
                    </div>
                </div>
            </div>
        )
    }

    return (
        <div className="bg-white rounded-lg shadow overflow-hidden">
            <div className="px-6 py-4 border-b border-gray-200">
                <h3 className="text-lg font-medium text-gray-900">
                    Stages ({forms?.length || 0})
                </h3>
            </div>
            <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                    <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Stage</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Section</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Statut formulaire</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Statut section</th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Mis à jour</th>
                    </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                    {forms.map(({ form, form_section }) => (
                        <tr key={form.id} className="hover:bg-gray-50">
                            <td className="px-6 py-4 whitespace-nowrap">
                                <div className="text-sm font-medium text-gray-900">Stage #{form.stage_id}</div>
                                <div className="text-sm text-gray-500">Étudiant #{form.student_id}</div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                                <div className="text-sm text-gray-900">{form_section.section_type}</div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                                    <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusStyle(form.status)}`}>
                                        {form.status}
                                    </span>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                                    <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusStyle(form_section.status)}`}>
                                        {form_section.status}
                                    </span>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {formatDate(form.updated_at)}
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>
        </div>
    )
}
