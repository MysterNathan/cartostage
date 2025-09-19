// app/student/page.tsx
'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { authApi } from '@/lib/authApi'
import type { StudentData } from '@/types/student'
import type { StudentApplication } from '@/types/student'
import ApplicationModal from '@/components/student/ApplicationModal'

export default function MyStudentPage() {
    const router = useRouter()
    const [data, setData] = useState<StudentData | null>(null)
    const [applications, setApplications] = useState<StudentApplication[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')

    // États pour les modals
    const [editingApplication, setEditingApplication] = useState<StudentApplication | null>(null)
    const [isApplicationModalOpen, setIsApplicationModalOpen] = useState(false)
    const [isApplicationNew, setIsApplicationNew] = useState(false)

    useEffect(() => {
        if (!authApi.isAuthenticated()) {
            router.push('/login')
            return
        }
    }, [router])

    useEffect(() => {
        loadData()
    }, [])

    const loadData = async () => {
        try {
            const studentData = await getMyStudentData()
            setData(studentData)
        } catch (error) {
            console.error('Erreur:', error)
            setError('Erreur lors du chargement des données')
        } finally {
            setLoading(false)
        }
    }


    const handleLogout = () => {
        authApi.logout()
        router.push('/login')
    }


    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-500 mx-auto mb-4"></div>
                    <p className="text-gray-600">Chargement...</p>
                </div>
            </div>
        )
    }

    if (error) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="text-center">
                    <p className="text-red-600 mb-4">{error}</p>
                    <button
                        onClick={loadData}
                        className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
                    >
                        Réessayer
                    </button>
                </div>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-gray-100">
            {/* Header */}
            <div className="bg-white shadow-sm border-b">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex justify-between items-center h-16">
                        <div>
                            <h1 className="text-xl font-semibold text-gray-900">
                                Mon Espace Étudiant
                            </h1>
                            {data && (
                                <p className="text-sm text-gray-600">
                                    {data.prenom} {data.nom} - {data.filiere}
                                </p>
                            )}
                        </div>
                        <div className="flex items-center gap-4">
                            <a
                                href="/"
                                className="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
                            >
                                Retour à la carte
                            </a>
                            <button
                                onClick={handleLogout}
                                className="px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors"
                            >
                                Déconnexion
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            {/* Contenu principal */}
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <StudentStats
                    applications={applications}
                />

                <ApplicationsList
                    applications={applications}
                    loading={loading}
                    onEdit={handleEditApplication}
                    onDelete={handleApplicationDelete}
                    onAdd={handleAddApplication}
                />
            </div>

            {/* Modal */}
            {isApplicationModalOpen && (
                <ApplicationModal
                    application={editingApplication}
                    isOpen={isApplicationModalOpen}
                    onClose={() => setIsApplicationModalOpen(false)}
                    onSuccess={handleApplicationSuccess}
                    onDelete={handleApplicationDelete}
                    isNew={isApplicationNew}
                />
            )}
        </div>
    )
}
