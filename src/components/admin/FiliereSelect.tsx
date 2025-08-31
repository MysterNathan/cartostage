'use client'
import { useMemo, useState } from 'react'
import { useFilieres } from '@/hooks/useFilieres'
import type { Filiere } from '@/types/filiere'


interface Props {
    value?: string // code sélectionné
    onChange: (code: string) => void
    allowCreate?: boolean
    label?: string
    required?: boolean
    className?: string
}


export default function FiliereSelect({ value, onChange, allowCreate = true, label = 'Filière', required, className }: Props) {
    const { filieres, create, loading, error } = useFilieres()
    const [showAdd, setShowAdd] = useState(false)
    const [form, setForm] = useState({ code: '', label: '', color: '#3B82F6' })


    const options = useMemo(() => filieres, [filieres])


    return (
        <div className={className}>
            <label className="block text-sm font-medium text-gray-700 mb-2">{label}{required ? ' *' : ''}</label>


            <div className="flex gap-2">
                <select
                    value={value ?? ''}
                    onChange={(e) => onChange(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    required={required}
                    disabled={loading}
                >
                    <option value="" disabled>Choisir une filière</option>
                    {options.map((f: Filiere) => (
                        <option key={f.id} value={f.code}>
                            {f.label} ({f.code})
                        </option>
                    ))}
                </select>


                {allowCreate && (
                    <button
                        type="button"
                        onClick={() => setShowAdd(v => !v)}
                        className="px-3 py-2 text-sm bg-gray-100 hover:bg-gray-200 rounded-md"
                        aria-expanded={showAdd}
                    >
                        + Ajouter
                    </button>
                )}
            </div>


            {error && <p className="text-xs text-red-600 mt-1">{error}</p>}


            {showAdd && (
                <div className="mt-3 p-3 border rounded-md bg-gray-50 space-y-2">
                    <div className="grid grid-cols-3 gap-2">
                        <input
                            type="text" placeholder="Code (ex: CCST)"
                            value={form.code}
                            onChange={(e) => setForm({ ...form, code: e.target.value.toUpperCase() })}
                            className="px-2 py-1 border rounded"
                        />
                        <input
                            type="text" placeholder="Libellé"
                            value={form.label}
                            onChange={(e) => setForm({ ...form, label: e.target.value })}
                            className="px-2 py-1 border rounded col-span-2"
                        />
                    </div>
                    <div className="flex items-center gap-2">
                        <label className="text-sm text-gray-600">Couleur :</label>
                        <input type="color" value={form.color} onChange={(e) => setForm({ ...form, color: e.target.value })} />
                    </div>
                    <div className="flex gap-2">
                        <button
                            type="button"
                            onClick={async () => {
                                if (!form.code || !form.label) return
                                try {
                                    const created = await create(form)
                                    if (created) {
                                        onChange(created.code)
                                        setShowAdd(false)
                                        setForm({ code: '', label: '', color: '#3B82F6' })
                                    }
                                } catch (e) {
                                    alert('Impossible de créer la filière')
                                }
                            }}
                            className="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700 text-sm"
                        >
                            Enregistrer
                        </button>
                        <button
                            type="button"
                            onClick={() => setShowAdd(false)}
                            className="px-3 py-1 bg-gray-200 rounded hover:bg-gray-300 text-sm"
                        >
                            Annuler
                        </button>
                    </div>
                </div>
            )}
        </div>
    )
}