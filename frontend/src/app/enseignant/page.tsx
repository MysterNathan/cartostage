// app/enterprise/page.tsx
'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { authApi } from '@/lib/api/authApi'
import {getStudentsUser, getTutorUser} from '@/lib/api/enterpriseApi'
import type { Tutor } from '@/types/tutor'
import TutorsList from '@/components/misc/TutorsList'
import TutorModal from '@/components/enterprise/TutorModal'
import StudentsList from '@/components/misc/StudentsList'
import {Student, Tutors} from "@/types/user";
import FormSectionsList from "@/components/form/FormSectionList";
import FormSectionModal from "@/components/form/FormSectionModal";
import {FormResponse, FormSection} from "@/types/form";
import {getForm} from "@/lib/api/stageApi";

export default function MyEnterprisePage() {
    const router = useRouter()
    const [tutors, setTutors] = useState<Tutors[]>([])
    const [students, setStudents] = useState<Student[]>([])
    const [loading, setLoading] = useState(true)

    // États pour les modals
    const [editingTutor, setEditingTutor] = useState<Tutor | null>(null)
    const [isTutorModalOpen, setIsTutorModalOpen] = useState(false)
    const [isTutorNew, setIsTutorNew] = useState(false)
    const [formLoading, setFormLoading] = useState(false)
    const [formResponse, setFormResponse] = useState<FormResponse | null>(null)

    useEffect(() => {
        if (!authApi.isAuthenticated()) {
            router.push('/login')
            return
        }
        if (!authApi.isTeacher()){
            router.push('/')
            return
        }
        loadData()
    }, [router])

    const loadData = async () => {
        try {
            setLoading(true)
            const tutorsDatas = await getTutorUser()
            const studentDatas = await getStudentsUser()
            setTutors(tutorsDatas)
            setStudents(studentDatas)
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

    return (
        <div className="min-h-screen bg-gray-100">
            {/* Contenu principal */}
            <button
                onClick={handleOpenForm}
                disabled={formLoading}
                className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
            >
                {formLoading ? "Chargement..." : "Mon formulaire"}
            </button>

            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <TutorsList initialTutors={tutors} loading={false} />
            </div>

            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <StudentsList students={students} />
            </div>

            {formResponse && (
                <FormSectionModal
                    formResponse={formResponse}
                    onClose={() => setFormResponse(null)}
                />
            )}


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
        </div>
    )
}
