// lib/api/enterpriseApi.ts
import type { EnterpriseData, EnterpriseStats } from '@/types/enterprise'
import {fetchApi} from "@/lib/api/apiHelper";

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
