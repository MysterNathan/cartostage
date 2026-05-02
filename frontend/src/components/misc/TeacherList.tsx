
// components/enterprise-service/TeachersList.tsx
'use client'

import { useState, useEffect } from 'react'
import type { User } from '@/types/user'
import { getUserDisplayName } from '@/lib/userUtils'
import {getTeacher} from "@/lib/api/tutorApi";

interface TeachersListProps {
    initialTeachers?: User[]
    autoFetch?: boolean
}

export default function TeachersList({ initialTeachers, autoFetch = true }: TeachersListProps) {
    const [teachers, setTeachers] = useState<User[]>(initialTeachers || [])
    const [loading, setLoading] = useState(autoFetch && !initialTeachers)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        if (autoFetch && !initialTeachers) {
            loadTeachers()
        }
    }, [autoFetch, initialTeachers])

    const loadTeachers = async () => {
        try {
            setLoading(true)
            setError(null)
            const data = await getTeacher()
            setTeachers(data || [])
        } catch (error) {
            console.error('Erreur lors du chargement des enseignants:', error)
            setError(error instanceof Error ? error.message : 'Erreur de chargement')
        } finally {
            setLoading(false)
        }
    }

    const getUserInitials = (user: User): string => {
        const first = user.first_name?.[0]?.toUpperCase() || ''
        const last = user.last_name?.[0]?.toUpperCase() || ''
        return first + last || user.username[0]?.toUpperCase() || '?'
    }

    const getUserPhone = (user: User): string | null => {
        return user.profile?.phone || null
    }

    const getUserPoste = (user: User): string => {
        return user.profile?.poste || 'Non défini'
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
                            onClick={loadTeachers}
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
                        <p className="text-gray-600">Chargement des enseignants...</p>
                    </div>
                </div>
            </div>
        )
    }

    return (
        <div className="bg-white rounded-lg shadow overflow-hidden">
            <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
                <h3 className="text-lg font-medium text-gray-900">
                    Enseignants ({teachers?.length || 0})
                </h3>
            </div>
            <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                    <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Enseignant
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
                    {teachers.map((teacher) => (
                        <tr key={teacher.id} className="hover:bg-gray-50">
                            <td className="px-6 py-4 whitespace-nowrap">
                                <div className="flex items-center">
                                    <div className="flex-shrink-0 h-10 w-10">
                                        <div className="h-10 w-10 rounded-full bg-purple-500 flex items-center justify-center">
                                                <span className="text-sm font-medium text-white">
                                                    {getUserInitials(teacher)}
                                                </span>
                                        </div>
                                    </div>
                                    <div className="ml-4">
                                        <div className="text-sm font-medium text-gray-900">
                                            {getUserDisplayName(teacher)}
                                        </div>
                                        <div className="text-sm text-gray-500">
                                            @{teacher.username}
                                        </div>
                                    </div>
                                </div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                                <div className="text-sm text-gray-900">{teacher.email}</div>
                                {getUserPhone(teacher) && (
                                    <div className="text-sm text-gray-500">{getUserPhone(teacher)}</div>
                                )}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                                    <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-purple-100 text-purple-800">
                                        {getUserPoste(teacher)}
                                    </span>
                                {teacher.profile?.departement && (
                                    <div className="text-xs text-gray-500 mt-1">
                                        {teacher.profile.departement}
                                    </div>
                                )}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                                    <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                        teacher.is_active && (teacher.profile?.is_active !== false)
                                            ? 'bg-green-100 text-green-800'
                                            : 'bg-red-100 text-red-800'
                                    }`}>
                                        {teacher.is_active && (teacher.profile?.is_active !== false) ? 'Actif' : 'Inactif'}
                                    </span>
                                {!teacher.email_verified && (
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
    )
}
