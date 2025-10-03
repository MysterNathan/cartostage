'use client'
import {useEffect, useState} from "react";
import {getEnterpriseStats} from "@/lib/enterpriseApi";
import authApi from "@/lib/authApi";
import {router} from "next/client";
import EnterpriseStats from "@/components/enterprise/EnterpriseStats";
import TutorsList from "@/components/student/TutorsList";
import {Tutor} from "@/types/tutors";

export default function StudentPage() {
    const [loading, setLoading] = useState(true)
    const [tutors, setTutors] = useState<Tutor[]>([])

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
            //const studentDatas = await getStudentDatas()
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
    return(
        <>
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
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <TutorsList
                    tutors={tutors}
                    loading={false}
                />
            </div>
        </>
    )

}