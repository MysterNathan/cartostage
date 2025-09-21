// app/enterprise/page.tsx
'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { authApi } from '@/lib/authApi'
import { getMyEnterpriseData } from '@/lib/enterpriseApi'
import type { EnterpriseWithStats } from '@/types/enterprise'
import type { EnterpriseData } from '@/types/enterprise'
import type { Tutor } from '@/types/tutor'
import EnterpriseStats from '@/components/enterprise/EnterpriseStats'
import TutorsList from '@/components/enterprise/TutorsList'
import TutorModal from '@/components/enterprise/TutorModal'
import StudentsList from '@/components/enterprise/StudentsList'

export default function MyEnterprisePage() {
    const router = useRouter()
    const [data, setData] = useState<EnterpriseData | null>(null)

    const [enterprise, setEnterprise] = useState<EnterpriseWithStats | null>(null)
    const [tutors, setTutors] = useState<Tutor[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')

    // États pour les modals
    const [editingTutor, setEditingTutor] = useState<Tutor | null>(null)
    const [isTutorModalOpen, setIsTutorModalOpen] = useState(false)
    const [isTutorNew, setIsTutorNew] = useState(false)

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
            const enterpriseData = await getMyEnterpriseData()
            setData(enterpriseData)
        } catch (error) {
            console.error('Erreur:', error)
        } finally {
            setLoading(false)
        }
    }

    const handleTutorAdded = (newTutor: Tutor) => {
        setData(prevData => ({
            ...prevData!,
            tutors: [...(prevData?.tutors || []), newTutor]
        }))
    }

    const handleTutorUpdated = (updatedTutor: Tutor) => {
        setData(prevData => ({
            ...prevData!,
            tutors: prevData!.tutors.map(tutor =>
                tutor.id === updatedTutor.id ? updatedTutor : tutor
            )
        }))
    }

    const handleLogout = () => {
        authApi.logout()
        router.push('/login')
    }


    // Gestionnaires pour les tuteurs
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
            setTutors([...tutors, tutor])
        } else {
            setTutors(tutors.map(t => t.id === tutor.id ? tutor : t))
        }
    }

    const handleTutorDelete = (tutorId: number) => {
        setTutors(tutors.filter(t => t.id !== tutorId))
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

    // Calcul des statistiques
    const safeTutors = tutors || []
    const activeTutors = safeTutors.length

    return (
        <div className="min-h-screen bg-gray-100">
            {/* Header */}
            <div className="bg-white shadow-sm border-b">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex justify-between items-center h-16">
                        <div>
                            <h1 className="text-xl font-semibold text-gray-900">
                                {enterprise?.nom || 'Mon Entreprise'}
                            </h1>
                            {enterprise?.adresse && (
                                <p className="text-sm text-gray-600">{enterprise.adresse}</p>
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
                <EnterpriseStats
                    activeTutors={activeTutors}
                />

                <TutorsList
                    tutors={data?.tutors || []}
                    loading={loading}
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
            {/* Nouvelle section pour les élèves */}
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <StudentsList />
            </div>
        </div>
    )
}
