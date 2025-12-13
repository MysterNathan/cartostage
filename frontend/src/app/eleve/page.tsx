'use client'
import {useEffect, useState} from "react";
import {getEnterpriseStats} from "@/lib/api/enterpriseApi";
import authApi from "@/lib/api/authApi";
import {router} from "next/client";
import EnterpriseStats from "@/components/enterprise/EnterpriseStats";
import TutorsList from "@/components/misc/TutorsList";
import {Tutor} from "@/types/tutors";
import TeachersList from "@/components/misc/TeacherList";
import {useRouter} from "next/navigation";

export default function StudentPage() {
    const [loading, setLoading] = useState(true)
    const [tutors, setTutors] = useState<Tutor[]>([])
    const router = useRouter()

    useEffect(() => {
        if (!authApi.isAuthenticated()) {
            router.push('/login')
            return
        }
        if (!authApi.isStudent()){
            router.push('/')
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
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <TutorsList
                    tutors={tutors}
                    loading={false}
                />
            </div>
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <TeachersList
                tutors={tutors}
                loading={false}
            />
        </div>
        </>
    )

}