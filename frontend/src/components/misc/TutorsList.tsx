// components/enterprise-service/TutorsList.tsx
'use client'

import { useState, useEffect } from 'react'
import { formatTutorData } from '@/lib/enterpriseUtils'
import type { User, } from '@/types/user'
import {getUserDisplayName} from '@/lib/userUtils'
import {getTutors} from "@/lib/api/tutorApi";

interface TutorsListProps {
    // Props optionnelles pour override les données
    initialTutors?: User[]
    autoFetch?: boolean
}

export default function TutorsList({ initialTutors, autoFetch = true }: TutorsListProps) {
    const [tutors, setTutors] = useState<User[]>(initialTutors || [])
    const [loading, setLoading] = useState(autoFetch && !initialTutors)
    const [error, setError] = useState<string | null>(null)

    // Fetch des tuteurs au montage si autoFetch est activé
    useEffect(() => {
        if (autoFetch && !initialTutors) {
            loadTutors()
        }
    }, [autoFetch, initialTutors])

    const loadTutors = async () => {
        try {
            setLoading(true)
            setError(null)
            const data = await getTutors()
            setTutors(data || [])
        } catch (error) {
            console.error('Erreur lors du chargement des tuteurs:', error)
            setError(error instanceof Error ? error.message : 'Erreur de chargement')
        } finally {
            setLoading(false)
        }
    }


    // Fonction helper pour obtenir les initiales
    const getUserInitials = (user: User): string => {
        const first = user.first_name?.[0]?.toUpperCase() || ''
        const last = user.last_name?.[0]?.toUpperCase() || ''
        return first + last || user.username[0]?.toUpperCase() || '?'
    }

    // Fonction helper pour obtenir le téléphone
    const getUserPhone = (user: User): string | null => {
        return user.profile?.phone || null
    }

    // Fonction helper pour obtenir le poste
    const getUserPoste = (user: User): string => {
        return user.profile?.poste || 'Non défini'
    }

    // Afficher l'erreur si il y en a une
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
                            onClick={loadTutors}
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
                        <p className="text-gray-600">Chargement des tuteurs...</p>
                    </div>
                </div>
            </div>
        )
    }

    return (
        <>
            <div className="bg-white rounded-lg shadow overflow-hidden">
                <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
                    <h3 className="text-lg font-medium text-gray-900">
                        Tuteurs ({tutors?.length || 0})
                    </h3>
                </div>
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-gray-200">
                            <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Tuteur
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Contact
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Poste
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Statut
                                </th>
                            </tr>
                            </thead>
                            <tbody className="bg-white divide-y divide-gray-200">
                            {tutors.map((tutor) => (
                                <tr key={tutor.id} className="hover:bg-gray-50">
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <div className="flex items-center">
                                            <div className="flex-shrink-0 h-10 w-10">
                                                <div className="h-10 w-10 rounded-full bg-blue-500 flex items-center justify-center">
                                                    <span className="text-sm font-medium text-white">
                                                        {getUserInitials(tutor)}
                                                    </span>
                                                </div>
                                            </div>
                                            <div className="ml-4">
                                                <div className="text-sm font-medium text-gray-900">
                                                    {getUserDisplayName(tutor)}
                                                </div>
                                                <div className="text-sm text-gray-500">
                                                    @{tutor.username}
                                                </div>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <div className="text-sm text-gray-900">{tutor.email}</div>
                                        {getUserPhone(tutor) && (
                                            <div className="text-sm text-gray-500">{getUserPhone(tutor)}</div>
                                        )}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                                            {getUserPoste(tutor)}
                                        </span>
                                        {tutor.profile?.departement && (
                                            <div className="text-xs text-gray-500 mt-1">
                                                {tutor.profile.departement}
                                            </div>
                                        )}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                            tutor.is_active && (tutor.profile?.is_active !== false)
                                                ? 'bg-green-100 text-green-800'
                                                : 'bg-red-100 text-red-800'
                                        }`}>
                                            {tutor.is_active && (tutor.profile?.is_active !== false) ? 'Actif' : 'Inactif'}
                                        </span>
                                        {!tutor.email_verified && (
                                            <div className="text-xs text-orange-600 mt-1">
                                                Email non vérifié
                                            </div>
                                        )}
                                    </td>
                                </tr>
                            ))}
                            </tbody>
                        </table>
                    </div>
            </div>
        </>
    )
}
