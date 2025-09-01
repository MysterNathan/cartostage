// src/app/api/stages/route.ts
import { NextResponse } from 'next/server'
import path from 'path'
import fs from 'fs/promises'

export const dynamic = 'force-dynamic'
const filePath = path.join(process.cwd(), 'public', 'data', 'stages.json')

type Stage = any

async function readFileSafe(): Promise<{ raw: any; exists: boolean }> {
    try {
        const txt = await fs.readFile(filePath, 'utf8')
        return { raw: JSON.parse(txt), exists: true }
    } catch {
        return { raw: [], exists: false }
    }
}

function extractStages(raw: any): { stages: Stage[]; shape: 'array' | 'object' } {
    if (Array.isArray(raw)) return { stages: raw, shape: 'array' }
    if (raw && Array.isArray(raw.stages)) return { stages: raw.stages, shape: 'object' }
    return { stages: [], shape: 'array' }
}

async function writeBack(shape: 'array' | 'object', stages: Stage[]) {
    const dataToWrite = shape === 'array' ? stages : { stages }
    const tmp = `${filePath}.tmp`
    await fs.mkdir(path.dirname(filePath), { recursive: true })
    await fs.writeFile(tmp, JSON.stringify(dataToWrite, null, 2), 'utf8')
    await fs.rename(tmp, filePath)
}

export async function GET() {
    const { raw } = await readFileSafe()
    const { stages } = extractStages(raw)
    return NextResponse.json({ stages }, { headers: { 'Cache-Control': 'no-store' } })
}

/**
 * Body attendu: { stages: Stage[] }
 * Écrit en préservant la forme d'origine du fichier (array pur vs {stages: []})
 */
export async function POST(req: Request) {
    try {
        const body = await req.json()
        const incoming = body?.stages
        if (!Array.isArray(incoming)) {
            return NextResponse.json({ error: 'Payload invalide: stages[] requis' }, { status: 400 })
        }

        const { raw } = await readFileSafe()
        const { shape } = extractStages(raw)
        await writeBack(shape, incoming)
        return NextResponse.json({ ok: true })
    } catch {
        return NextResponse.json({ error: 'Écriture échouée' }, { status: 500 })
    }
}
