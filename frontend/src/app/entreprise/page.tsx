// app/enterprise/page.tsx
'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { authApi } from '@/lib/authApi'
import { getEnterpriseStats } from '@/lib/enterpriseApi'
import type { EnterpriseStats } from '@/types/enterprise'
import type { Tutor } from '@/types/tutor'
import EnterpriseStats from '@/components/enterprise/EnterpriseStats'
import TutorsList from '@/components/enterprise/TutorsList'
import TutorModal from '@/components/enterprise/TutorModal'
import StudentsList from '@/components/enterprise/StudentsList'

export default function MyEnterprisePage() {
    const router = useRouter()
    const [stats, setStats] = useState<EnterpriseStats | null>(null)
    const [tutors, setTutors] = useState<Tutor[]>([])
    const [loading, setLoading] = useState(true)

    // États pour les modals
    const [editingTutor, setEditingTutor] = useState<Tutor | null>(null)
    const [isTutorModalOpen, setIsTutorModalOpen] = useState(false)
    const [isTutorNew, setIsTutorNew] = useState(false)

    useEffect(() => {
        if (!authApi.isAuthenticated()) {
            router.push('/login')
            return
        }
        loadData()
    }, [router])

    const loadData = async () => {
        try {
            setLoading(true)
            const enterpriseStats = await getEnterpriseStats()
            setStats(enterpriseStats)
        } catch (error) {
            console.error('Erreur lors du chargement des données:', error)
        } finally {
            setLoading(false)
        }
    }

    const handleTutorAdded = (newTutor: Tutor) => {
        setTutors(prevTutors => [...prevTutors, newTutor])
        loadData() // Recharger les stats
    }

    const handleTutorUpdated = (updatedTutor: Tutor) => {
        setTutors(prevTutors =>
            prevTutors.map(tutor =>
                tutor.id === updatedTutor.id ? updatedTutor : tutor
            )
        )
    }

    const handleTutorDelete = (tutorId: number) => {
        setTutors(prevTutors => prevTutors.filter(t => t.id !== tutorId))
        loadData() // Recharger les stats
    }

    const handleLogout = () => {
        authApi.logout()
        router.push('/login')
    }

    const handleAddTutor = () => {
        setEditingTutor(null)
        setIsTutorNew(true)
        setIsTutorModalOpen(true)
    }

    const handleEditTutor = (tutor: Tutor) => {
        setEditingTutor(tutor)
        setIsTutorNew(false)
        setIsTutorModalOpen(true)
    }

    const handleTutorSuccess = (tutor: Tutor) => {
        if (isTutorNew) {
            handleTutorAdded(tutor)
        } else {
            handleTutorUpdated(tutor)
        }
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

    if (!stats) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="text-center">
                    <p className="text-gray-600">Aucune donnée disponible</p>
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
                                Mon Entreprise
                            </h1>
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
                <EnterpriseStats
                    activeTutors={stats.tutors}
                    totalStages={stats.stages}
                    totalStudents={stats.students}
                />

                <TutorsList
                    tutors={tutors}
                    loading={false}
                    onTutorAdded={handleTutorAdded}
                    onTutorUpdated={handleTutorUpdated}
                />
            </div>

            {/* Modal */}
            {isTutorModalOpen && (
                <TutorModal
                    tutor={editingTutor}
                    isOpen={isTutorModalOpen}
                    onClose={() => setIsTutorModalOpen(false)}
                    onSuccess={handleTutorSuccess}
                    onDelete={handleTutorDelete}
                    isNew={isTutorNew}
                />
            )}

            {/* Section élèves */}
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <StudentsList />
            </div>
        </div>
    )
}
