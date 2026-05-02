// lib/errorUtils.ts

/**
 * Gère les erreurs d'API et retourne un message utilisateur approprié
 */
export const handleApiError = (error: unknown): string => {
    console.error('Erreur API:', error);

    if (error instanceof Error) {
        // Erreurs d'authentification
        if (error.message.includes('401') || error.message.includes('Unauthorized')) {
            return 'Session expirée. Veuillez vous reconnecter.';
        }

        // Erreurs de validation
        if (error.message.includes('400') || error.message.includes('Bad Request')) {
            return 'Données invalides. Veuillez vérifier les informations saisies.';
        }

        // Erreurs de serveur
        if (error.message.includes('500') || error.message.includes('Internal Server Error')) {
            return 'Erreur serveur. Veuillez réessayer plus tard.';
        }

        // Erreurs réseau
        if (error.message.includes('Failed to fetch') || error.message.includes('Network')) {
            return 'Problème de connexion. Vérifiez votre connexion internet.';
        }

        // Conflits (ex: stage-service déjà existant)
        if (error.message.includes('409') || error.message.includes('Conflict')) {
            return 'Conflit détecté. Cette donnée existe peut-être déjà.';
        }

        // Retourner le message d'erreur s'il est lisible
        if (error.message.length < 100 && !error.message.includes('Error:')) {
            return error.message;
        }
    }

    return 'Une erreur inattendue s\'est produite. Veuillez réessayer.';
};

/**
 * Détermine si une erreur nécessite une reconnexion
 */
export const isAuthError = (error: unknown): boolean => {
    if (error instanceof Error) {
        return error.message.includes('401') ||
            error.message.includes('Unauthorized') ||
            error.message.includes('Token expired');
    }
    return false;
};
