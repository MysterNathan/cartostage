// app/student-service/page.tsx
'use client'

import { useEffect, useState } from "react"
import authApi from "@/lib/api/authApi"
import TutorsList from "@/components/misc/TutorsList"
import TeachersList from "@/components/misc/TeacherList"
import { useRouter } from "next/navigation"
import FormSectionModal from "@/components/form/FormSectionModal"
import {FormResponse} from "@/types/form";
import FormsList from "@/components/form/FormList";
import {getForm} from "@/lib/api/stageApi";


export default function StudentPage() {
    const [formResponse, setFormResponse] = useState<FormResponse | null>(null)
    const [formLoading, setFormLoading] = useState(false)

    const router = useRouter()


    useEffect(() => {
        if (!authApi.isAuthenticated()) {
            router.push('/login')
            return
        }
        if (!authApi.isStudent()) {
            router.push('/')
            return
        }
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
    return (
        <div className="min-h-screen bg-gray-50">
            {/* Header de page */}
            <div className="bg-white border-b border-gray-200 shadow-sm">
                <div className="max-w-[1920px] mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
                    <div>
                        <h1 className="text-xl font-semibold text-gray-900">Espace Élève</h1>
                        <p className="text-sm text-gray-500">Gérez vos contacts et consultez les stages disponibles</p>
                    </div>
                    <div className="flex items-center gap-3">
                        <button
                            onClick={() => router.push('/carte')}
                            className="flex items-center gap-2 px-4 py-2 bg-green-500 text-white text-sm font-medium rounded-lg hover:bg-green-600 shadow-sm transition-colors"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
                            </svg>
                            Accéder à la carte
                        </button>
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
            </div>

            {/* Contenu principal */}
            <div className="max-w-[1920px] mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <div className="grid grid-cols-1 xl:grid-cols-2 gap-8 mb-8">
                    <div className="w-full">
                        <TutorsList/>
                    </div>
                    <div className="w-full">
                        <TeachersList/>
                    </div>
                </div>
                <div>
                    <FormsList/>
                </div>
            </div>

            {/* Modals */}
            {formResponse && (
                <FormSectionModal
                    formResponse={formResponse}
                    onClose={() => setFormResponse(null)}
                />
            )}
        </div>
    )
}
