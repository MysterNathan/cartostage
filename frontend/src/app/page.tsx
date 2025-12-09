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

            {/* Body - Contenu principal avec diagonale */}
            <main className="flex-grow relative overflow-hidden">
                <div className="absolute inset-0">
                    {/* Zone gauche (en haut à gauche) - Background bleu */}
                    <div
                        className="absolute top-0 left-0 w-full h-full bg-blue-50"
                        style={{
                            clipPath: 'polygon(0 0, 100% 100%, 0 100%)'
                        }}
                    />

                    {/* Zone droite (en bas à droite) - Background blanc */}
                    <div
                        className="absolute top-0 right-0 w-full h-full bg-white"
                        style={{
                            clipPath: 'polygon(100% 0, 100% 100%, 0 100%)'
                        }}
                    />

                    {/* Image Array (en haut à gauche) - Au-dessus des zones clippées */}
                    <div className="absolute top-0 left-0 w-3/5 h-3/5 flex items-start justify-start p-6 z-10">
                        <div className="relative w-full h-full">
                            <Image
                                src="/carte/array.png"
                                alt="Tableau des stages"
                                fill
                                className="object-contain object-left-top"
                                priority
                            />
                        </div>
                    </div>

                    {/* Image Map (en bas à droite) - Au-dessus des zones clippées */}
                    <div className="absolute bottom-0 right-0 w-3/5 h-3/5 flex items-end justify-end p-6 z-10">
                        <div className="relative w-full h-full">
                            <Image
                                src="/carte/map.png"
                                alt="Carte des stages"
                                fill
                                className="object-contain object-right-bottom"
                                priority
                            />
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
