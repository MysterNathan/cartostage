'use client'
import { useState, useEffect } from 'react'
import { getFilieres, addFiliere, updateFiliere, deleteFiliere } from '@/lib/filiereApi'
import type { Filiere } from '@/types/filiere'

interface Props {
    isOpen: boolean
    onClose: () => void
}

export default function FilieresManager({ isOpen, onClose }: Props) {
    const [filieres, setFilieres] = useState<Filiere[]>([])
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [draft, setDraft] = useState<Partial<Filiere> | null>(null)
    const [createForm, setCreateForm] = useState({ code: '', label: '', color: '#3B82F6' })

    // Charger les filières
    const loadFilieres = async () => {
        try {
            setLoading(true)
            setError(null)
            const response = await getFilieres()
            setFilieres(response.filieres)
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Erreur lors du chargement')
        } finally {
            setLoading(false)
        }
    }

    // Charger au montage et quand la modal s'ouvre
    useEffect(() => {
        if (isOpen) {
            loadFilieres()
        }
    }, [isOpen])

    // Créer une filière
    const create = async (data: { code: string; label: string; color: string }) => {
        try {
            setLoading(true)
            setError(null)
            const newFiliere = await addFiliere(data)
            setFilieres(prev => [...prev, newFiliere])
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Erreur lors de la création')
            throw err
        } finally {
            setLoading(false)
        }
    }

    // Modifier une filière
    const update = async (data: { id: number; code?: string; label?: string; color?: string }) => {
        try {
            setLoading(true)
            setError(null)
            const updatedFiliere = await updateFiliere(data.id, data)
            setFilieres(prev => prev.map(f => f.id === data.id ? updatedFiliere : f))
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Erreur lors de la modification')
            throw err
        } finally {
            setLoading(false)
        }
    }

    // Supprimer une filière
    const remove = async (id: number) => {
        try {
            setLoading(true)
            setError(null)
            await deleteFiliere(id)
            setFilieres(prev => prev.filter(f => f.id !== id))
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Erreur lors de la suppression')
            throw err
        } finally {
            setLoading(false)
        }
    }

    if (!isOpen) return null

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4" onClick={(e) => e.target === e.currentTarget && onClose()}>
            <div className="absolute inset-0 bg-black/50" />
            <div className="relative bg-white rounded-xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-hidden">

                {/* Header */}
                <div className="flex items-center justify-between p-6 border-b border-gray-200 bg-slate-50">
                    <h2 className="text-xl font-bold text-gray-900">Gestion des filières</h2>
                    <button
                        onClick={onClose}
                        className="p-2 hover:bg-gray-200 rounded-lg transition-colors text-gray-500 hover:text-gray-700"
                    >
                        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                <div className="p-6 space-y-6 max-h-[calc(90vh-140px)] overflow-y-auto">

                    {/* Message d'erreur */}
                    {error && (
                        <div className="p-4 bg-red-50 border border-red-200 rounded-lg">
                            <div className="text-sm text-red-800 font-medium">{error}</div>
                        </div>
                    )}

                    {/* Formulaire d'ajout */}
                    <div className="p-5 border border-gray-200 rounded-xl bg-blue-50">
                        <h3 className="font-semibold text-gray-900 mb-4 flex items-center">
                            <svg className="w-5 h-5 mr-2 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                            </svg>
                            Ajouter une nouvelle filière
                        </h3>
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Code</label>
                                <input
                                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white text-gray-900"
                                    placeholder="Ex: INFO"
                                    value={createForm.code}
                                    onChange={(e) => setCreateForm({ ...createForm, code: e.target.value.toUpperCase() })}
                                />
                            </div>
                            <div className="md:col-span-2">
                                <label className="block text-sm font-medium text-gray-700 mb-1">Libellé</label>
                                <input
                                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white text-gray-900"
                                    placeholder="Ex: Informatique"
                                    value={createForm.label}
                                    onChange={(e) => setCreateForm({ ...createForm, label: e.target.value })}
                                />
                            </div>
                        </div>
                        <div className="flex items-center justify-between">
                            <div className="flex items-center gap-3">
                                <label className="text-sm font-medium text-gray-700">Couleur :</label>
                                <div className="flex items-center gap-2">
                                    <input
                                        type="color"
                                        value={createForm.color}
                                        onChange={(e) => setCreateForm({ ...createForm, color: e.target.value })}
                                        className="w-10 h-10 border border-gray-300 rounded-lg cursor-pointer"
                                    />
                                    <span className="text-sm text-gray-600">{createForm.color}</span>
                                </div>
                            </div>
                            <button
                                disabled={!createForm.code || !createForm.label || loading}
                                onClick={async () => {
                                    await create(createForm)
                                    setCreateForm({ code: '', label: '', color: '#3B82F6' })
                                }}
                                className="px-6 py-2 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-medium disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                            >
                                {loading ? 'Création...' : 'Créer la filière'}
                            </button>
                        </div>
                    </div>

                    {/* Liste des filières */}
                    <div className="border border-gray-200 rounded-xl overflow-hidden bg-white">
                        <div className="bg-gray-50 px-6 py-3 border-b border-gray-200">
                            <h3 className="font-semibold text-gray-900">Liste des filières ({filieres.length})</h3>
                        </div>

                        {loading && filieres.length === 0 ? (
                            <div className="p-8 text-center">
                                <div className="animate-spin w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full mx-auto mb-3"></div>
                                <p className="text-gray-600">Chargement des filières...</p>
                            </div>
                        ) : filieres.length === 0 ? (
                            <div className="p-8 text-center">
                                <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-3">
                                    <svg className="w-8 h-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                                    </svg>
                                </div>
                                <p className="text-gray-600 font-medium">Aucune filière trouvée</p>
                                <p className="text-gray-500 text-sm">Commencez par ajouter votre première filière</p>
                            </div>
                        ) : (
                            <div className="overflow-x-auto">
                                <table className="w-full">
                                    <thead className="bg-gray-50">
                                    <tr>
                                        <th className="text-left p-4 font-semibold text-gray-900 text-sm">Code</th>
                                        <th className="text-left p-4 font-semibold text-gray-900 text-sm">Libellé</th>
                                        <th className="text-left p-4 font-semibold text-gray-900 text-sm">Couleur</th>
                                        <th className="text-right p-4 font-semibold text-gray-900 text-sm">Actions</th>
                                    </tr>
                                    </thead>
                                    <tbody className="divide-y divide-gray-200">
                                    {filieres.map((f) => (
                                        <tr key={f.id} className="hover:bg-gray-50 transition-colors">
                                            <td className="p-4">
                                                {draft?.id === f.id ? (
                                                    <input
                                                        className="w-24 px-3 py-1 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white text-gray-900"
                                                        value={draft.code ?? ''}
                                                        onChange={(e) => setDraft(d => ({ ...d!, code: e.target.value.toUpperCase() }))}
                                                    />
                                                ) : (
                                                    <span className="font-mono font-medium text-gray-900 bg-gray-100 px-2 py-1 rounded text-sm">
                                                            {f.code}
                                                        </span>
                                                )}
                                            </td>
                                            <td className="p-4">
                                                {draft?.id === f.id ? (
                                                    <input
                                                        className="w-full px-3 py-1 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white text-gray-900"
                                                        value={draft.label ?? ''}
                                                        onChange={(e) => setDraft(d => ({ ...d!, label: e.target.value }))}
                                                    />
                                                ) : (
                                                    <span className="text-gray-900 font-medium">{f.label}</span>
                                                )}
                                            </td>
                                            <td className="p-4">
                                                {draft?.id === f.id ? (
                                                    <div className="flex items-center gap-2">
                                                        <input
                                                            type="color"
                                                            value={draft.color ?? '#3B82F6'}
                                                            onChange={(e) => setDraft(d => ({ ...d!, color: e.target.value }))}
                                                            className="w-8 h-8 border border-gray-300 rounded cursor-pointer"
                                                        />
                                                        <span className="text-sm text-gray-600">{draft.color}</span>
                                                    </div>
                                                ) : (
                                                    <div className="flex items-center gap-3">
                                                        <div
                                                            className="w-6 h-6 rounded-full border border-gray-300"
                                                            style={{ backgroundColor: f.color }}
                                                        />
                                                        <span className="text-sm font-mono text-gray-600 bg-gray-50 px-2 py-1 rounded">
                                                                {f.color}
                                                            </span>
                                                    </div>
                                                )}
                                            </td>
                                            <td className="p-4 text-right">
                                                {draft?.id === f.id ? (
                                                    <div className="flex gap-2 justify-end">
                                                        <button
                                                            onClick={async () => {
                                                                await update({
                                                                    id: draft.id!,
                                                                    code: draft.code,
                                                                    label: draft.label,
                                                                    color: draft.color
                                                                })
                                                                setDraft(null)
                                                            }}
                                                            className="px-3 py-1 text-white bg-green-600 hover:bg-green-700 rounded-md text-sm font-medium transition-colors"
                                                        >
                                                            Enregistrer
                                                        </button>
                                                        <button
                                                            onClick={() => setDraft(null)}
                                                            className="px-3 py-1 bg-gray-200 hover:bg-gray-300 text-gray-700 rounded-md text-sm font-medium transition-colors"
                                                        >
                                                            Annuler
                                                        </button>
                                                    </div>
                                                ) : (
                                                    <div className="flex gap-3 justify-end">
                                                        <button
                                                            onClick={() => setDraft({ ...f })}
                                                            className="text-blue-600 hover:text-blue-800 font-medium text-sm transition-colors"
                                                        >
                                                            Modifier
                                                        </button>
                                                        <button
                                                            onClick={async () => {
                                                                if (confirm(`Êtes-vous sûr de vouloir supprimer la filière "${f.label}" ?`)) {
                                                                    await remove(f.id)
                                                                }
                                                            }}
                                                            className="text-red-600 hover:text-red-800 font-medium text-sm transition-colors"
                                                        >
                                                            Supprimer
                                                        </button>
                                                    </div>
                                                )}
                                            </td>
                                        </tr>
                                    ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </div>
                </div>

                {/* Footer */}
                <div className="flex justify-end gap-3 p-6 border-t border-gray-200 bg-gray-50">
                    <button
                        onClick={onClose}
                        className="px-6 py-2 bg-white border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 font-medium transition-colors"
                    >
                        Fermer
                    </button>
                </div>
            </div>
        </div>
    )
}
