// components/enterprise-service/TutorModal.tsx
import { useState, useEffect } from 'react'
import type { Tutor } from '@/types/tutor'

interface TutorModalProps {
    tutor: Tutor | null
    isOpen: boolean
    onClose: () => void
    onSuccess: (tutor: Tutor) => void
    onDelete?: (tutorId: number) => void
    isNew: boolean
}

export default function TutorModal({
                                       tutor,
                                       isOpen,
                                       onClose,
                                       onSuccess,
                                       onDelete,
                                       isNew
                                   }: TutorModalProps) {
    const [formData, setFormData] = useState({
        nom: '',
        prenom: '',
        email: '',
        telephone: '',
        poste: ''
    })
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState('')

    useEffect(() => {
        if (isOpen) {
            if (tutor && !isNew) {
                setFormData({
                    nom: tutor.nom || '',
                    prenom: tutor.prenom || '',
                    email: tutor.email || '',
                    telephone: tutor.telephone || '',
                    poste: tutor.poste || ''
                })
            } else {
                setFormData({
                    nom: '',
                    prenom: '',
                    email: '',
                    telephone: '',
                    poste: ''
                })
            }
            setError('')
        }
    }, [tutor, isOpen, isNew])

    if (!isOpen) return null

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setLoading(true)
        setError('')

        try {
            // Cette logique sera implémentée avec les vraies API calls
            const mockTutor: Tutor = {
                id: tutor?.id || Date.now(),
                ...formData,
                enterprise_id: 1, // Sera récupéré du contexte
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString()
            }

            onSuccess(mockTutor)
            onClose()
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Une erreur est survenue')
        } finally {
            setLoading(false)
        }
    }

    const handleDelete = async () => {
        if (!tutor || !onDelete) return

        if (confirm('Êtes-vous sûr de vouloir supprimer ce tuteur ?')) {
            try {
                setLoading(true)
                onDelete(tutor.id)
                onClose()
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Erreur lors de la suppression')
            } finally {
                setLoading(false)
            }
        }
    }

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md max-h-[90vh] overflow-y-auto">
                <div className="flex justify-between items-center mb-6">
                    <h2 className="text-xl font-semibold">
                        {isNew ? 'Ajouter un tuteur' : 'Modifier le tuteur'}
                    </h2>
                    <button
                        onClick={onClose}
                        className="text-gray-400 hover:text-gray-600"
                    >
                        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                {error && (
                    <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                Nom *
                            </label>
                            <input
                                type="text"
                                required
                                value={formData.nom}
                                onChange={(e) => setFormData({ ...formData, nom: e.target.value })}
                                className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                Prénom *
                            </label>
                            <input
                                type="text"
                                required
                                value={formData.prenom}
                                onChange={(e) => setFormData({ ...formData, prenom: e.target.value })}
                                className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Email *
                        </label>
                        <input
                            type="email"
                            required
                            value={formData.email}
                            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                            className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Téléphone
                        </label>
                        <input
                            type="tel"
                            value={formData.telephone}
                            onChange={(e) => setFormData({ ...formData, telephone: e.target.value })}
                            className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Poste *
                        </label>
                        <input
                            type="text"
                            required
                            value={formData.poste}
                            onChange={(e) => setFormData({ ...formData, poste: e.target.value })}
                            className="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            placeholder="ex: Responsable informatique, Chef de projet..."
                        />
                    </div>

                    <div className="flex justify-between pt-4">
                        <div>
                            {!isNew && onDelete && (
                                <button
                                    type="button"
                                    onClick={handleDelete}
                                    disabled={loading}
                                    className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 disabled:opacity-50"
                                >
                                    Supprimer
                                </button>
                            )}
                        </div>
                        <div className="flex gap-2">
                            <button
                                type="button"
                                onClick={onClose}
                                disabled={loading}
                                className="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 disabled:opacity-50"
                            >
                                Annuler
                            </button>
                            <button
                                type="submit"
                                disabled={loading}
                                className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50"
                            >
                                {loading ? 'Sauvegarde...' : (isNew ? 'Ajouter' : 'Modifier')}
                            </button>
                        </div>
                    </div>
                </form>
            </div>
        </div>
    )
}
