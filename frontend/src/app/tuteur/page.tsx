// app/enterprise/page.tsx
'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { authApi } from '@/lib/api/authApi'
import {getEnterpriseStats, getStudentsUser} from '@/lib/api/enterpriseApi'
import type { EnterpriseStats } from '@/types/enterprise'
import type { Tutor, Student } from '@/types/tutor'
import EnterpriseStats from '@/components/enterprise/EnterpriseStats'
import TutorsList from '@/components/misc/TutorsList'
import TutorModal from '@/components/enterprise/TutorModal'
import StudentsList from '@/components/misc/StudentsList'
import {getForm} from "@/lib/api/stageApi";
import {FormResponse} from "@/types/form";
import FormSectionModal from "@/components/form/FormSectionModal";

export default function MyEnterprisePage() {
    const router = useRouter()
    const [stats, setStats] = useState<EnterpriseStats | null>(null)
    const [tutors, setTutors] = useState<Tutor[]>([])
    const [students, setStudents] = useState<Student[]>([])
    const [loading, setLoading] = useState(true)
    const [formLoading, setFormLoading] = useState(false)

    const [editingTutor, setEditingTutor] = useState<Tutor | null>(null)
    const [isTutorModalOpen, setIsTutorModalOpen] = useState(false)
    const [isTutorNew, setIsTutorNew] = useState(false)
    const [formResponse, setFormResponse] = useState<FormResponse | null>(null)

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

    const handleOpenForm = async () => {
        setFormLoading(true)
        try {
            const data = await getForm()
            console.log("datas:",data)
            setFormResponse(data)
        } catch (error) {
            console.error("Erreur lors du chargement du formulaire", error)
        }
        setFormLoading(false)
    }

    const loadData = async () => {
        try {
            setLoading(true)
            const enterpriseStats = await getEnterpriseStats()
            setStats(enterpriseStats)
            const studentsUsers = await getStudentsUser()
            setStudents(studentsUsers)
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
            <div>
                <button
                    onClick={handleOpenForm}
                    disabled={formLoading}
                    className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
                >
                    {formLoading ? "Chargement..." : "Mon formulaire"}
                </button>
            </div>
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
                <StudentsList
                    students={students}
                />
            </div>

            {formResponse && (
                <FormSectionModal
                    formResponse={formResponse}
                    onClose={() => setFormResponse(null)}
                />
            )}
        </div>
    )
}
