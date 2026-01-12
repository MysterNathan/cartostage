// app/student/page.tsx
'use client'

import { useEffect, useState } from "react"
import authApi from "@/lib/api/authApi"
import TutorsList from "@/components/misc/TutorsList"
import { Tutor } from "@/types/user"
import TeachersList from "@/components/misc/TeacherList"
import { useRouter } from "next/navigation"
import StageMapView from "@/components/student/MapModal"
import { getStages } from "@/lib/api/stageApi"

export default function StudentPage() {
    const [tutors, setTutors] = useState<Tutor[]>([])
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

    const handleMap = () => {
        router.push('/carte')
    }
    
    return (
        <div className="min-h-screen bg-gray-50">
            <div className="max-w-[1920px] mx-auto px-4 sm:px-6 lg:px-8 py-8">
                {/* En-tête avec bouton carte */}
                <div className="mb-8 grid grid-cols-1 md:grid-cols-2 gap-6 items-start">
                    <div>
                        <h1 className="text-3xl font-bold text-gray-900">Espace Élève</h1>
                        <p className="mt-2 text-gray-600">Gérez vos contacts et consultez les stages disponibles</p>
                    </div>

                    <div className="flex justify-start md:justify-end">
                        <button
                            onClick={handleMap}
                            className="px-6 py-3 bg-green-400 text-white font-semibold rounded-lg shadow-md hover:bg-green-500 transition-colors duration-200 cursor-pointer"
                        >
                            Accéder à la carte
                        </button>
                    </div>
                </div>

                {/* Tuteurs et Enseignants côte à côte - Plus larges */}
                <div className="grid grid-cols-1 xl:grid-cols-2 gap-8 mb-8">
                    {/* Tuteurs */}
                    <div className="w-full">
                        <TutorsList
                            tutors={tutors}
                            loading={false}
                        />
                    </div>

                    {/* Enseignants */}
                    <div className="w-full">
                        <TeachersList
                            tutors={tutors}
                            loading={false}
                        />
                    </div>
                </div>
            </div>
        </div>
    )
}
