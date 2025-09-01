// src/app/api/filieres/route.ts
import { NextResponse } from 'next/server'
import path from 'path'
import fs from 'fs/promises'

export const dynamic = 'force-dynamic'
const filieresPath = path.join(process.cwd(), 'public', 'data', 'filieres.json')
const stagesPath = path.join(process.cwd(), 'public', 'data', 'stages.json')

type Filiere = { id: number; code: string; label: string; color: string }

async function readJSON<T>(file: string, fallback: T): Promise<T> {
    try { return JSON.parse(await fs.readFile(file, 'utf8')) as T } catch { return fallback }
}
async function writeJSON(file: string, data: unknown) {
    const tmp = `${file}.tmp`
    await fs.mkdir(path.dirname(file), { recursive: true })
    await fs.writeFile(tmp, JSON.stringify(data, null, 2), 'utf8')
    await fs.rename(tmp, file)
}

export async function GET() {
    const filieres = await readJSON<Filiere[]>(filieresPath, [])
    return NextResponse.json({ filieres }, { headers: { 'Cache-Control': 'no-store' } })
}

export async function POST(req: Request) {
    const body = await req.json()
    let { code, label, color } = body ?? {}
    if (!code || !label) return NextResponse.json({ error: 'code et label requis' }, { status: 400 })

    code = String(code).toUpperCase().trim()
    color = (color ?? '#3B82F6').trim()

    const filieres = await readJSON<Filiere[]>(filieresPath, [])
    const existing = filieres.find(f => f.code === code)
    if (existing) {
        existing.label = label
        existing.color = color
        await writeJSON(filieresPath, filieres)
        return NextResponse.json({ filiere: existing }, { status: 201 })
    }

    const id = filieres.length ? Math.max(...filieres.map(f => f.id)) + 1 : 1
    const created: Filiere = { id, code, label, color }
    filieres.push(created)
    await writeJSON(filieresPath, filieres)
    return NextResponse.json({ filiere: created }, { status: 201 })
}

export async function PUT(req: Request) {
    const { id, code, label, color } = await req.json()
    if (!id) return NextResponse.json({ error: 'id requis' }, { status: 400 })

    const filieres = await readJSON<Filiere[]>(filieresPath, [])
    const f = filieres.find(x => x.id === Number(id))
    if (!f) return NextResponse.json({ error: 'introuvable' }, { status: 404 })

    if (code) f.code = String(code).toUpperCase().trim()
    if (label) f.label = String(label).trim()
    if (color) f.color = String(color).trim()

    await writeJSON(filieresPath, filieres)
    return NextResponse.json({ filiere: f })
}

export async function DELETE(req: Request) {
    try {
        const { id } = await req.json()
        if (!id) return NextResponse.json({ error: 'id requis' }, { status: 400 })

        const filieres = await readJSON<Filiere[]>(filieresPath, [])
        const idx = filieres.findIndex(x => x.id === Number(id))
        if (idx === -1) return NextResponse.json({ error: 'introuvable' }, { status: 404 })

        // ⚠️ Normaliser stages.json (peut être [] ou { stages: [] })
        const stagesRaw = await readJSON<any>(stagesPath, [])
        const stagesArr = Array.isArray(stagesRaw)
            ? stagesRaw
            : (Array.isArray(stagesRaw?.stages) ? stagesRaw.stages : [])

        const code = filieres[idx].code
        if (stagesArr.some((s: any) => s?.filiere === code)) {
            return NextResponse.json({ error: 'Des stages utilisent cette filière' }, { status: 409 })
        }

        filieres.splice(idx, 1)
        await writeJSON(filieresPath, filieres)
        return NextResponse.json({ ok: true })
    } catch (e) {
        return NextResponse.json({ error: 'Suppression échouée' }, { status: 500 })
    }
}