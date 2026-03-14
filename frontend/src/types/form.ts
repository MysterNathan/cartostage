export interface Form {
    id: number
    stage_id: number
    student_id: number
    status: string
    content?: Record<string, unknown>
    created_at: string
    updated_at: string
}

export interface FormSection {
    id: number
    form_id: number
    section_type: string
    user_id: number
    status: string
    created_at: string
    updated_at: string
}


export interface FormResponse {
    form: Form
    form_section: FormSection
}

export interface FormResponses {
    form_responses: FormResponse[]
}