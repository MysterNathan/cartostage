// lib/stageUtils.ts
import { Stage } from '@/types/stage';

/**
 * Nettoie et valide les données d'un stage
 */
export const sanitizeStage = (stage: Partial<Stage>): Partial<Stage> => {
    return {
        ...stage,
        // Nettoyer les chaînes de caractères
        entreprise: stage.entreprise?.trim() || '',
        poste: stage.poste?.trim() || '',
        adresse: stage.adresse?.trim() || '',
        commune: stage.commune?.trim() || '',
        sector: stage.sector?.trim() || '',
        filiere: stage.filiere?.trim() || '',
        period: stage.period?.trim() || '',

        // Valider et nettoyer les nombres
        lat: typeof stage.lat === 'number' ? stage.lat : parseFloat(String(stage.lat)) || 0,
        lng: typeof stage.lng === 'number' ? stage.lng : parseFloat(String(stage.lng)) || 0,
        capacity_total: Math.max(0, parseInt(String(stage.capacity_total)) || 1),
        capacity_filled: Math.max(0, parseInt(String(stage.capacity_filled)) || 0),
        placesDisponibles: Math.max(0, parseInt(String(stage.placesDisponibles)) || 0),

        // Garder l'ID s'il existe
        ...(stage.id && { id: stage.id })
    };
};

/**
 * Calcule les places disponibles
 */
export const calculateAvailablePlaces = (total: number, filled: number): number => {
    return Math.max(0, total - filled);
};

/**
 * Valide les données d'un stage
 */
export const validateStage = (stage: Partial<Stage>): string | null => {
    if (!stage.entreprise?.trim()) {
        return 'Le nom de l\'enterprise est requis';
    }

    if (!stage.poste?.trim()) {
        return 'Le poste est requis';
    }

    if (!stage.adresse?.trim()) {
        return 'L\'adresse est requise';
    }

    if (!stage.filiere?.trim()) {
        return 'La filière est requise';
    }

    const total = stage.capacity_total || 0;
    const filled = stage.capacity_filled || 0;

    if (total < 1) {
        return 'La capacité totale doit être d\'au moins 1';
    }

    if (filled > total) {
        return 'Le nombre de places occupées ne peut pas dépasser la capacité totale';
    }

    return null;
};
