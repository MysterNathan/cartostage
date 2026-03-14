// app/student/page.tsx
'use client'

import { useEffect, useState } from "react"
import authApi from "@/lib/api/authApi"
import TutorsList from "@/components/misc/TutorsList"
import TeachersList from "@/components/misc/TeacherList"
import { useRouter } from "next/navigation"
import FormSectionModal from "@/components/form/FormSectionModal"
import {FormResponse} from "@/types/form";
import FormsList from "@/components/form/FormList";


export default function StudentPage() {
    const [formResponse, setFormResponse] = useState<FormResponse | null>(null)
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

    return (
        <>
            <div className="min-h-screen bg-gray-50">
                <div className="max-w-[1920px] mx-auto px-4 sm:px-6 lg:px-8 py-8">
                    <div className="mb-8 grid grid-cols-1 md:grid-cols-2 gap-6 items-start">
                        <div>
                            <h1 className="text-3xl font-bold text-gray-900">Espace Élève</h1>
                            <p className="mt-2 text-gray-600">Gérez vos contacts et consultez les stages disponibles</p>
                        </div>

                        <div className="flex justify-start md:justify-end gap-4">
                            <button
                                onClick={() => router.push('/carte')}
                                className="px-6 py-3 bg-green-400 text-white font-semibold rounded-lg shadow-md hover:bg-green-500 transition-colors duration-200 cursor-pointer"
                            >
                                Accéder à la carte
                            </button>

                            {formResponse && (
                            <FormSectionModal
                                formResponse={formResponse}
                                onClose={() => setFormResponse(null)}
                            />
                        )}

                        </div>
                    </div>

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
            </div>
        </>
    )
}
