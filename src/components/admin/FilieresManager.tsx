'use client'
import { useState } from 'react'
import { useFilieres } from '@/hooks/useFilieres'
import type { Filiere } from '@/types/filiere'

interface Props { isOpen: boolean; onClose: () => void }

export default function FilieresManager({ isOpen, onClose }: Props) {
    const { filieres, create, update, remove, loading, error } = useFilieres()
    const [draft, setDraft] = useState<Partial<Filiere> | null>(null)
    const [createForm, setCreateForm] = useState({ code: '', label: '', color: '#3B82F6' })
    if (!isOpen) return null

    return (
        <div className="fixed inset-0 z-50" onClick={(e) => e.target === e.currentTarget && onClose()}>
            <div className="absolute inset-0 backdrop-blur-sm bg-black/10" />
            <div className="relative mx-auto mt-10 bg-white border border-gray-200 rounded-lg shadow-2xl max-w-3xl w-[95%] overflow-hidden">
                <div className="flex items-center justify-between p-4 border-b">
                    <h2 className="text-lg font-semibold">Gestion des filières</h2>
                    <button onClick={onClose} className="p-2 hover:bg-gray-100 rounded">✕</button>
                </div>

                <div className="p-4 space-y-6 max-h-[70vh] overflow-y-auto">
                    {error && <div className="text-sm text-red-600">{error}</div>}

                    {/* Ajout */}
                    <div className="p-3 border rounded-lg bg-gray-50">
                        <h3 className="font-medium mb-2">Ajouter une filière</h3>
                        <div className="grid grid-cols-3 gap-2">
                            <input className="px-2 py-1 border rounded" placeholder="Code"
                                   value={createForm.code}
                                   onChange={(e) => setCreateForm({ ...createForm, code: e.target.value.toUpperCase() })}/>
                            <input className="px-2 py-1 border rounded col-span-2" placeholder="Libellé"
                                   value={createForm.label}
                                   onChange={(e) => setCreateForm({ ...createForm, label: e.target.value })}/>
                        </div>
                        <div className="flex items-center gap-2 mt-2">
                            <label className="text-sm text-gray-600">Couleur :</label>
                            <input type="color" value={createForm.color}
                                   onChange={(e) => setCreateForm({ ...createForm, color: e.target.value })}/>
                            <button
                                disabled={!createForm.code || !createForm.label || loading}
                                onClick={async () => { await create(createForm); setCreateForm({ code: '', label: '', color: '#3B82F6' }) }}
                                className="ml-auto px-3 py-1 rounded bg-blue-600 text-white text-sm disabled:opacity-50">
                                Créer
                            </button>
                        </div>
                    </div>

                    {/* Liste */}
                    <div className="border rounded-lg overflow-hidden">
                        <table className="w-full text-sm">
                            <thead className="bg-gray-50">
                            <tr>
                                <th className="text-left p-2">Code</th>
                                <th className="text-left p-2">Libellé</th>
                                <th className="text-left p-2">Couleur</th>
                                <th className="text-right p-2">Actions</th>
                            </tr>
                            </thead>
                            <tbody>
                            {filieres.map((f) => (
                                <tr key={f.id} className="border-t">
                                    <td className="p-2">
                                        {draft?.id === f.id
                                            ? <input className="px-2 py-1 border rounded w-28"
                                                     value={draft.code ?? ''} onChange={(e) => setDraft(d => ({ ...d!, code: e.target.value.toUpperCase() }))}/>
                                            : <span className="font-medium">{f.code}</span>}
                                    </td>
                                    <td className="p-2">
                                        {draft?.id === f.id
                                            ? <input className="px-2 py-1 border rounded w-full"
                                                     value={draft.label ?? ''} onChange={(e) => setDraft(d => ({ ...d!, label: e.target.value }))}/>
                                            : <span>{f.label}</span>}
                                    </td>
                                    <td className="p-2">
                                        {draft?.id === f.id
                                            ? <input type="color" value={draft.color ?? '#3B82F6'}
                                                     onChange={(e) => setDraft(d => ({ ...d!, color: e.target.value }))}/>
                                            : <span className="inline-flex items-center gap-2">
                            <span className="inline-block w-4 h-4 rounded" style={{ backgroundColor: f.color }} />
                            <span className="text-xs text-gray-500">{f.color}</span>
                          </span>}
                                    </td>
                                    <td className="p-2 text-right">
                                        {draft?.id === f.id
                                            ? <div className="inline-flex gap-2">
                                                <button onClick={async () => { await update({ id: draft.id!, code: draft.code, label: draft.label, color: draft.color }); setDraft(null) }}
                                                        className="px-2 py-1 text-white bg-blue-600 rounded">Enregistrer</button>
                                                <button onClick={() => setDraft(null)} className="px-2 py-1 bg-gray-200 rounded">Annuler</button>
                                            </div>
                                            : <div className="inline-flex gap-3">
                                                <button onClick={() => setDraft({ ...f })} className="text-blue-600">Modifier</button>
                                                <button onClick={async () => { if (confirm('Supprimer cette filière ?')) await remove(f.id) }}
                                                        className="text-red-600">Supprimer</button>
                                            </div>}
                                    </td>
                                </tr>
                            ))}
                            {filieres.length === 0 && (
                                <tr><td colSpan={4} className="p-4 text-center text-gray-500">Aucune filière</td></tr>
                            )}
                            </tbody>
                        </table>
                    </div>
                </div>

                <div className="flex justify-end gap-2 p-4 border-t">
                    <button onClick={onClose} className="px-4 py-2 bg-gray-100 rounded hover:bg-gray-200">Fermer</button>
                </div>
            </div>
        </div>
    )
}
