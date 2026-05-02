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

export interface Skill {
    name: string
    level: number
}

export interface Task {
    name: string
    summary: string
    skills: Skill[]
}

export interface FormSectionContentStudent {
    company_integration: number
    understanding_tasks: number
    autonomy_level: number
    skills_acquired: Skill[]
    positive_point: string
    difficulties_encountered: string
    comment: string
}

export interface FormSectionContentTeacher {
    attendance: number
    report_quality: number
    oral_presentation: number
    professional_conduct: number
    objectives_achieved: boolean
    grade: number
    recommandations: string
}

export interface FormSectionContentTutor {
    technical_skills: number
    work_quality: number
    punctuality: number
    team_integration: number
    autonomy: number
    tasks_completed: Task[]
    comment: string
}

export type FormSectionContent =
    | FormSectionContentStudent
    | FormSectionContentTeacher
    | FormSectionContentTutor
    | null

export interface FormFormSection {
    form: Form
    form_section: FormSection[]
    form_section_content: FormSectionContent
}

export type FormResponse = FormFormSection[]
