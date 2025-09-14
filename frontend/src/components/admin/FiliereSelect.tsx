'use client'
import { useState, useEffect, useMemo } from 'react'
import { useFilieres } from '@/hooks/useFilieres'
import type { Filiere } from '@/types/filiere'

interface Props {
    value?: string
    onChange: (code: string) => void
    allowCreate?: boolean
    label?: string
    required?: boolean
    className?: string
    disabled?: boolean
}

export default function FiliereSelect({
                                          value,
                                          onChange,
                                          allowCreate = true,
                                          label = 'Filière',
                                          required,
                                          className,
                                          disabled = false
                                      }: Props) {
    // Utilisation du hook personnalisé
    const { filieres, loading, error: hookError, loadFilieres, create } = useFilieres()

    // États locaux pour le formulaire d'ajout
    const [localError, setLocalError] = useState('')
    const [showAdd, setShowAdd] = useState(false)
    const [form, setForm] = useState({ code: '', label: '', color: '#3B82F6' })
    const [isCreating, setIsCreating] = useState(false)

    // Charger les filières au montage
    useEffect(() => {
        loadFilieres()
    }, [])

    const options = useMemo(() => {
        return Array.isArray(filieres) && filieres.length > 0 ? filieres : []
    }, [filieres])

    // Créer une nouvelle filière
    const handleCreate = async () => {
        if (!form.code.trim() || !form.label.trim()) {
            setLocalError('Le code et le libellé sont obligatoires')
            return
        }

        // Vérifier si le code existe déjà
        if (filieres.some(f => f.code === form.code.trim().toUpperCase())) {
            setLocalError('Ce code existe déjà')
            return
        }

        try {
            setIsCreating(true)
            setLocalError('')

            await create({
                code: form.code.trim().toUpperCase(),
                label: form.label.trim(),
                color: form.color
            })

            // Sélectionner la nouvelle filière
            onChange(form.code.trim().toUpperCase())

            // Réinitialiser le formulaire
            setShowAdd(false)
            setForm({ code: '', label: '', color: '#3B82F6' })

        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : 'Erreur lors de la création'
            setLocalError(errorMessage)
            console.error('Erreur création filière:', err)
        } finally {
            setIsCreating(false)
        }
    }

    // Réinitialiser l'erreur quand on change le formulaire
    useEffect(() => {
        if (localError && (form.code || form.label)) {
            setLocalError('')
        }
    }, [form.code, form.label, localError])

    const selectedFiliere = filieres.find(f => f.code === value)
    const displayError = localError || hookError

    return (
        <div className={className}>
            <label className="block text-sm font-medium text-gray-700 mb-2">
                {label}{required ? ' *' : ''}
            </label>

            <div className="flex gap-2">
                <div className="relative flex-1">
                    <select
                        value={value ?? ''}
                        onChange={(e) => onChange(e.target.value)}
                        disabled={disabled || loading}
                        required={required}
                        className="w-full pl-3 pr-10 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed appearance-none bg-white"
                    >
                        <option value="">
                            {loading ? 'Chargement...' : 'Sélectionner une filière'}
                        </option>
                        {options.map((f) => (
                            <option key={f.id} value={f.code}>
                                {f.label} ({f.code})
                            </option>
                        ))}
                    </select>

                    {/* Indicateur de couleur */}
                    {selectedFiliere && (
                        <div
                            className="absolute right-8 top-1/2 transform -translate-y-1/2 w-4 h-4 rounded-full border border-gray-300"
                            style={{ backgroundColor: selectedFiliere.color }}
                            title={`Couleur: ${selectedFiliere.color}`}
                        />
                    )}

                    {/* Icône de dropdown */}
                    <div className="absolute right-2 top-1/2 transform -translate-y-1/2 pointer-events-none">
                        <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                    </div>
                </div>

                {allowCreate && !disabled && (
                    <button
                        type="button"
                        onClick={() => setShowAdd(v => !v)}
                        className="px-3 py-2 text-sm bg-gray-100 hover:bg-gray-200 rounded-md transition-colors disabled:opacity-50"
                        aria-expanded={showAdd}
                        disabled={loading}
                    >
                        + Ajouter
                    </button>
                )}
            </div>

            {/* Affichage des erreurs */}
            {displayError && (
                <p className="text-xs text-red-600 mt-1 flex items-center gap-1">
                    <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                    </svg>
                    {displayError}
                </p>
            )}

            {/* Informations sur la filière sélectionnée */}
            {selectedFiliere && (
                <div className="mt-1 text-xs text-gray-500 flex items-center gap-2">
                    <span>Code: {selectedFiliere.code}</span>
                    <div
                        className="w-3 h-3 rounded-full border border-gray-300"
                        style={{ backgroundColor: selectedFiliere.color }}
                    />
                </div>
            )}

            {/* Formulaire d'ajout */}
            {showAdd && (
                <div className="mt-3 p-3 border rounded-md bg-gray-50 space-y-3">
                    <h4 className="text-sm font-medium text-gray-700">Nouvelle filière</h4>

                    <div className="grid grid-cols-3 gap-2">
                        <input
                            type="text"
                            placeholder="Code (ex: CCST)"
                            value={form.code}
                            onChange={(e) => setForm({ ...form, code: e.target.value.toUpperCase() })}
                            className="px-2 py-1 border rounded text-sm focus:outline-none focus:ring-1 focus:ring-blue-500"
                            disabled={isCreating}
                            maxLength={10}
                        />
                        <input
                            type="text"
                            placeholder="Libellé"
                            value={form.label}
                            onChange={(e) => setForm({ ...form, label: e.target.value })}
                            className="px-2 py-1 border rounded col-span-2 text-sm focus:outline-none focus:ring-1 focus:ring-blue-500"
                            disabled={isCreating}
                            maxLength={100}
                        />
                    </div>

                    <div className="flex items-center gap-2">
                        <label className="text-sm text-gray-600">Couleur :</label>
                        <input
                            type="color"
                            value={form.color}
                            onChange={(e) => setForm({ ...form, color: e.target.value })}
                            className="w-8 h-8 border rounded cursor-pointer disabled:opacity-50"
                            disabled={isCreating}
                        />
                        <span className="text-xs text-gray-500 font-mono">{form.color}</span>
                    </div>

                    <div className="flex gap-2">
                        <button
                            type="button"
                            onClick={handleCreate}
                            disabled={!form.code.trim() || !form.label.trim() || isCreating}
                            className="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700 text-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1"
                        >
                            {isCreating && (
                                <svg className="animate-spin w-3 h-3" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                    <path className="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                </svg>
                            )}
                            {isCreating ? 'Création...' : 'Enregistrer'}
                        </button>
                        <button
                            type="button"
                            onClick={() => {
                                setShowAdd(false)
                                setForm({ code: '', label: '', color: '#3B82F6' })
                                setLocalError('')
                            }}
                            disabled={isCreating}
                            className="px-3 py-1 bg-gray-200 rounded hover:bg-gray-300 text-sm transition-colors disabled:opacity-50"
                        >
                            Annuler
                        </button>
                    </div>
                </div>
            )}
        </div>
    )
}
