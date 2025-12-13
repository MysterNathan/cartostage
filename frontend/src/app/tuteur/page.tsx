// app/enterprise/page.tsx
'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { authApi } from '@/lib/api/authApi'
import { getEnterpriseStats } from '@/lib/api/enterpriseApi'
import type { EnterpriseStats } from '@/types/enterprise'
import type { Tutor } from '@/types/tutor'
import EnterpriseStats from '@/components/enterprise/EnterpriseStats'
import TutorsList from '@/components/misc/TutorsList'
import TutorModal from '@/components/enterprise/TutorModal'
import StudentsList from '@/components/misc/StudentsList'

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
        if (!authApi.isTutor()){
            router.push('/')
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
