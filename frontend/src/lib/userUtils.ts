// Fonctions utilitaires pour TypeScript
export const getUserFullName = (user: User): string => {
    if (user.first_name && user.last_name) {
        return `${user.first_name} ${user.last_name}`
    }
    if (user.first_name) return user.first_name
    if (user.last_name) return user.last_name
    return user.username
}

export const getUserDisplayName = (user: User): string => {
    const fullName = getUserFullName(user)
    return fullName !== user.username ? fullName : user.username
}

export const isUserActive = (user: User): boolean => {
    return user.is_active && (!user.profile || user.profile.is_active)
}

// Type guards
export const isTutor = (user: User): user is Tutor => {
    return user.role === 'tutor'
}

export const isStudent = (user: User): user is User => {
    return user.role === 'student'
}

export const isTeacher = (user: User): user is User => {
    return user.role === 'teacher'
}

export const isAdmin = (user: User): user is User => {
    return user.role === 'admin'
}

export const isEnterprise = (user: User): user is User => {
    return user.role === 'enterprise'
}
