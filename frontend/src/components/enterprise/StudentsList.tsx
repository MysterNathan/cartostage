// components/enterprise/StudentsList.tsx
'use client'

import { useState, useEffect } from 'react'
import { addStudent, updateStudent, deleteStudent, getStudents } from '@/lib/enterpriseApi'
import { formatStudentData } from '@/lib/enterpriseUtils'
import type { User } from '@/types/user'
import { getUserDisplayName } from '@/lib/userUtils'
import StudentStageModal from "@/components/enterprise/StudentStageModal";

interface StudentsListProps {
    // Props optionnelles pour override les données
    initialStudents?: User[]
    autoFetch?: boolean
}

export default function StudentsList({ initialStudents, autoFetch = true }: StudentsListProps) {
    const [students, setStudents] = useState<User[]>(initialStudents || [])
    const [loading, setLoading] = useState(autoFetch && !initialStudents)
    const [error, setError] = useState<string | null>(null)
    const [showForm, setShowForm] = useState(false)
    const [editingStudent, setEditingStudent] = useState<User | null>(null)
    const [showStageModal, setShowStageModal] = useState<User | null>(null)


    // Fetch des étudiants au montage si autoFetch est activé
    useEffect(() => {
        if (autoFetch && !initialStudents) {
            loadStudents()
        }
    }, [autoFetch, initialStudents])

    const loadStudents = async () => {
        try {
            setLoading(true)
            setError(null)
            const data = await getStudents()
            setStudents(data || [])
        } catch (error) {
            console.error('Erreur lors du chargement des élèves:', error)
            setError(error instanceof Error ? error.message : 'Erreur de chargement')
        } finally {
            setLoading(false)
        }
    }

    const handleAdd = () => {
        setEditingStudent(null)
        setShowForm(true)
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

    // Fonction helper pour obtenir la filière
    const getUserFiliere = (user: User): string => {
        return user.profile?.filiere || 'Non définie'
    }

    // Fonction helper pour obtenir l'établissement
    const getUserEtablissement = (user: User): string => {
        return user.profile?.etablissement || 'Non défini'
    }

    // Afficher l'erreur si présente
    if (error) {
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
                            onClick={loadStudents}
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
                        <p className="text-gray-600">Chargement des élèves...</p>
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
                        Élèves ({students?.length || 0})
                    </h3>
                </div>
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-gray-200">
                            <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Élève
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Contact
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Formation
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Statut
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Actions
                                </th>
                            </tr>
                            </thead>
                            <tbody className="bg-white divide-y divide-gray-200">
                            {students.map((student) => (
                                <tr key={student.id} className="hover:bg-gray-50">
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <div className="flex items-center">
                                            <div className="flex-shrink-0 h-10 w-10">
                                                <div className="h-10 w-10 rounded-full bg-green-500 flex items-center justify-center">
                                                    <span className="text-sm font-medium text-white">
                                                        {getUserInitials(student)}
                                                    </span>
                                                </div>
                                            </div>
                                            <div className="ml-4">
                                                <div className="text-sm font-medium text-gray-900">
                                                    {getUserDisplayName(student)}
                                                </div>
                                                <div className="text-sm text-gray-500">
                                                    @{student.username}
                                                </div>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <div className="text-sm text-gray-900">{student.email}</div>
                                        {getUserPhone(student) && (
                                            <div className="text-sm text-gray-500">{getUserPhone(student)}</div>
                                        )}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-purple-100 text-purple-800">
                                            {getUserFiliere(student)}
                                        </span>
                                        {getUserEtablissement(student) && (
                                            <div className="text-xs text-gray-500 mt-1">
                                                {getUserEtablissement(student)}
                                            </div>
                                        )}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                            student.is_active && (student.profile?.is_active !== false)
                                                ? 'bg-green-100 text-green-800'
                                                : 'bg-red-100 text-red-800'
                                        }`}>
                                            {student.is_active && (student.profile?.is_active !== false) ? 'Actif' : 'Inactif'}
                                        </span>
                                        {!student.email_verified && (
                                            <div className="text-xs text-orange-600 mt-1">
                                                Email non vérifié
                                            </div>
                                        )}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">(
                                        <div className="flex items-center space-x-3">
                                            <button
                                                onClick={() => setShowStageModal(student)}
                                                className="text-green-600 hover:text-green-900 font-medium flex items-center"
                                            >
                                                <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                                                </svg>
                                                Voir le stage
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            ))}
                            </tbody>
                        </table>
                    </div>
            </div>
            <StudentStageModal
                student={showStageModal}
                onClose={() => setShowStageModal(null)}
            />
        </>
    )
}
