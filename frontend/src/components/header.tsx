"use client";

import { useRouter, usePathname } from "next/navigation";
import authApi from "@/lib/api/authApi";
import { useEffect, useState } from "react";

// Configuration des rôles
const ROLE_CONFIG = {
    admin: {
        label: "Administration",
        route: "/admin"
    },
    student: {
        label: "Espace Étudiant",
        route: "/eleve"
    },
    teacher: {
        label: "Espace Enseignant",
        route: "/enseignant"
    },
    tutor: {
        label: "Espace Tuteur",
        route: "/tuteur"
    }
} as const;

type UserRole = keyof typeof ROLE_CONFIG;

export default function Header() {
    const router = useRouter();
    const pathname = usePathname(); // Détecte les changements de route
    const [loading, setLoading] = useState(true);
    const [role, setRole] = useState<string | null>(null);

    useEffect(() => {
        // Fonction pour récupérer le rôle
        const fetchRole = () => {
            const userRole = window.localStorage.getItem("role");
            setRole(userRole);
        };

        fetchRole();
        setLoading(false);
    }, [pathname]); // Se déclenche à chaque changement de page

    const handleHome = () => {
        router.push("/");
    };

    const handleLogin = () => {
        router.push("/login");
    };

    const handleAreaButton = () => {
        if (role && role in ROLE_CONFIG) {
            router.push(ROLE_CONFIG[role as UserRole].route);
        } else {
            router.push("/login");
        }
    };

    const getRoleButtonText = () => {
        if (role && role in ROLE_CONFIG) {
            return ROLE_CONFIG[role as UserRole].label;
        }
        return "Mon Espace";
    };

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-500 mx-auto mb-4"></div>
                    <p className="text-gray-600">Chargement...</p>
                </div>
            </div>
        );
    }

    return (
        <header className="bg-white shadow-sm border-b">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex justify-between items-center h-16">
                    <button
                        onClick={handleHome}
                        className="hover:scale-110 transition-transform"
                    >
                        <h1 className="text-xl font-semibold text-gray-900">
                            CartoStage
                        </h1>
                    </button>
                    <div className="flex items-center gap-4">
                        {authApi.isAuthenticated() ? (
                            <>
                                <button
                                    onClick={handleAreaButton}
                                    className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 hover:scale-105 transition-all"
                                >
                                    {getRoleButtonText()}
                                </button>
                                <button
                                    onClick={() => {
                                        authApi.logout();
                                        window.location.reload();
                                    }}
                                    className="px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 hover:scale-105 transition-all"
                                >
                                    Déconnexion
                                </button>
                            </>
                        ) : (
                            <button
                                onClick={handleLogin}
                                className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 hover:scale-105 transition-all"
                            >
                                Connexion
                            </button>
                        )}
                    </div>
                </div>
            </div>
        </header>
    );
}
