// lib/api/tutorApi.ts
import type { Tutor } from '@/types/tutor'
import {fetchApi} from "@/lib/api/apiHelper";

export async function getTutors(): Promise<Tutor[]> {
    return fetchApi<Tutor[]>('/api/enterprises/users')
}

export async function getTutorById(id: number): Promise<Tutor> {
    return fetchApi<Tutor>(`/api/enterprises/${id}`)
}

export async function getTeacher(): Promise<Tutor[]> {
    return fetchApi<Tutor[]>('/api/teacher/users')
}

export async function getStudents(): Promise<Tutor[]> {
    return fetchApi<Tutor[]>('/api/student/users')
}

export async function addTutor(
    tutorData: Omit<Tutor, 'id' | 'enterprise_id' | 'created_at' | 'updated_at'>
): Promise<Tutor> {
    return fetchApi<Tutor>('/api/enterprises', {
        method: 'POST',
        body: JSON.stringify(tutorData),
    })
}

export async function updateTutor(
    id: number,
    tutorData: Partial<Tutor>
): Promise<Tutor> {
    return fetchApi<Tutor>(`/api/enterprises/${id}`, {
        method: 'PUT',
        body: JSON.stringify(tutorData),
    })
}

export async function deleteTutor(id: number): Promise<void> {
    return fetchApi<void>(`/api/enterprises/${id}`, {
        method: 'DELETE',
    })
}
