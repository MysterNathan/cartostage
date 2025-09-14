// lib/filieres.ts
import type { Filiere } from "@/types"  // adapte l'import selon où est ton type

/**
 * Charge les filières depuis ton backend / API.
 * @returns Liste des filières
 * @throws Erreur si la requête échoue
 */
export async function fetchFilieres(): Promise<Filiere[]> {
    try {
        // 🔹 À remplacer par ton vrai appel API ou service
        const response = await fetch("/api/filieres");
        if (!response.ok) {
            throw new Error("Impossible de charger les filières");
        }

        const data = await response.json();

        return Array.isArray(data.filieres) ? data.filieres : [];
    } catch (err) {
        throw err instanceof Error
            ? err
            : new Error("Erreur inconnue lors du chargement des filières");
    }
}
