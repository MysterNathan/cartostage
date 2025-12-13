'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { authApi } from '@/lib/api/authApi'
import Image from 'next/image'

export default function HomePage() {
    const router = useRouter()
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        // Vérifier l'authentification au chargement
        setLoading(false)
    }, [])

    const handleLogin = () => {
        router.push('/login')
    }

    const handleAdmin = () => {
        if (authApi.isAdmin()) {
            router.push('/admin')
        } else {
            router.push('/login')
        }
    }

    // Affichage du loader pendant le chargement initial
    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-500 mx-auto mb-4"></div>
                    <p className="text-gray-600">Chargement...</p>
                </div>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-gray-100 flex flex-col">
            {/* Header */}
            <header className="bg-white shadow-sm border-b">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex justify-between items-center h-16">
                        <h1 className="text-xl font-semibold text-gray-900">Carte des stages</h1>
                        <div className="flex items-center gap-4">
                            {authApi.isAuthenticated() ? (
                                <>
                                    {authApi.isAdmin() && (
                                        <button
                                            onClick={handleAdmin}
                                            className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                                        >
                                            Administration
                                        </button>
                                    )}
                                    <button
                                        onClick={() => {
                                            authApi.logout()
                                            window.location.reload()
                                        }}
                                        className="px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors"
                                    >
                                        Déconnexion
                                    </button>
                                </>
                            ) : (
                                <button
                                    onClick={handleLogin}
                                    className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                                >
                                    Connexion
                                </button>
                            )}
                        </div>
                    </div>
                </div>
            </header>

            {/* Body - 3 cases pour les espaces */}
            <main className="flex-grow py-12 px-4 bg-white"
                  style={{
                      backgroundImage: `linear-gradient(#f3f4f6 1px, transparent 1px),
                      linear-gradient(90deg, #f3f4f6 1px, transparent 1px)`,
                      backgroundSize: '50px 50px'
                  }}
            >
                <div className="max-w-7xl mx-auto">
                    <h2 className="text-3xl font-bold text-center text-gray-900 mb-12">
                        Choisissez votre espace
                    </h2>

                    <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                        {/* Espace Élève */}
                        <div className="bg-white rounded-lg shadow-lg overflow-hidden flex flex-col transition-transform duration-300 hover:scale-105 cursor-pointer">
                            <div className="relative h-64 bg-gray-200">
                                <Image
                                    src="/carte/map.png"
                                    alt="Espace Élève"
                                    fill
                                    className="object-cover"
                                />
                            </div>
                            <div className="p-6 flex-grow flex flex-col">
                                <h3 className="text-2xl font-semibold text-gray-900 mb-2">
                                    Espace Élève
                                </h3>
                                <p className="text-gray-600 mb-4 flex-grow">
                                    Consultez la carte des stages disponibles
                                </p>
                                <button
                                    onClick={() => {
                                        // À compléter
                                    }}
                                    className="w-full px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                                >
                                    Accéder
                                </button>
                            </div>
                        </div>

                        {/* Espace Enseignant */}
                        <div className="bg-white rounded-lg shadow-lg overflow-hidden flex flex-col transition-transform duration-300 hover:scale-105 cursor-pointer">
                            <div className="relative h-64 bg-gray-200">
                                <Image
                                    src="/carte/array.png"
                                    alt="Espace Enseignant"
                                    fill
                                    className="object-cover"
                                />
                            </div>
                            <div className="p-6 flex-grow flex flex-col">
                                <h3 className="text-2xl font-semibold text-gray-900 mb-2">
                                    Espace Enseignant
                                </h3>
                                <p className="text-gray-600 mb-4 flex-grow">
                                    Gérez les stages et suivez vos élèves
                                </p>
                                <button
                                    onClick={() => {
                                        // À compléter
                                    }}
                                    className="w-full px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                                >
                                    Accéder
                                </button>
                            </div>
                        </div>

                        {/* Espace Tuteur */}
                        <div className="bg-white rounded-lg shadow-lg overflow-hidden flex flex-col transition-transform duration-300 hover:scale-105 cursor-pointer">
                            <div className="relative h-64 bg-gray-200 flex items-center justify-center">
                                <span className="text-gray-400 text-lg">Image à venir</span>
                            </div>
                            <div className="p-6 flex-grow flex flex-col">
                                <h3 className="text-2xl font-semibold text-gray-900 mb-2">
                                    Espace Tuteur
                                </h3>
                                <p className="text-gray-600 mb-4 flex-grow">
                                    Accédez aux informations de vos stagiaires
                                </p>
                                <button
                                    onClick={() => {
                                        // À compléter
                                    }}
                                    className="w-full px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                                >
                                    Accéder
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </main>

            {/* Footer */}
            <footer className="bg-white border-t">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
                    <div className="text-center text-gray-600 text-sm">
                        <p>© {new Date().getFullYear()} Carte des stages - Tous droits réservés</p>
                    </div>
                </div>
            </footer>
        </div>
    )
}
