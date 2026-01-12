// lib/api/enterpriseApi.ts
import type { EnterpriseData, EnterpriseStats } from '@/types/enterprise'
import {fetchApi} from "@/lib/api/apiHelper";
import {Student, Teacher, Tutor} from "@/types/user";

export async function getMyEnterpriseData(): Promise<EnterpriseData> {
    const data = await fetchApi<any>('/api/enterprises/me')

    return {
        enterprise: data.enterprise || null,
        stages: data.stages || [],
        tutors: data.tutors || [],
    }
}

export async function getEnterpriseStats(): Promise<EnterpriseStats> {
    return fetchApi<EnterpriseStats>('/api/enterprises/stats')
}

export async function getStudentsUser(): Promise<Student[]> {
    return fetchApi<Student[]>('/api/student/users')
}

export async function getTeacherUsers(): Promise<Teacher[]> {
    return fetchApi<Teacher[]>('/api/teacher/users')
}

export async function getTutorUser(): Promise<Tutor[]> {
    return fetchApi<Tutor[]>('/api/enterprises/users')
}