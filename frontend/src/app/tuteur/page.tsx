// app/enterprise-service/page.tsx
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
            {/* Header de page */}
            <div className="bg-white border-b border-gray-200 shadow-sm">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
                    <div>
                        <h1 className="text-xl font-semibold text-gray-900">Espace Entreprise</h1>
                        <p className="text-sm text-gray-500">Gérez vos tuteurs et vos élèves</p>
                    </div>
                    <button
                        onClick={handleOpenForm}
                        disabled={formLoading}
                        className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 shadow-sm transition-colors"
                    >
                        {formLoading ? (
                            <>
                                <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/>
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                                </svg>
                                Chargement...
                            </>
                        ) : (
                            <>
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                </svg>
                                Mes fiches de stage
                            </>
                        )}
                    </button>
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
                />
            </div>

            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <StudentsList students={students} />
            </div>

            {/* Modals */}
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
            {formResponse && (
                <FormSectionModal
                    formResponse={formResponse}
                    onClose={() => setFormResponse(null)}
                />
            )}
        </div>
    )

}
